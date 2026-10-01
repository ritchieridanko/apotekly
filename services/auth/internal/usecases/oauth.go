package usecases

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/ritchieridanko/apotekly/services/auth/internal/infra/oauth"
	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/contract/events/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/publisher"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"go.opentelemetry.io/otel"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OAuthUsecase interface {
	GetAuthURL(ctx context.Context, provider string) (url string, err *ce.Error)
	Callback(ctx context.Context, req *models.OAuthCallbackReq) (at *models.AuthToken, err *ce.Error)
}

type oAuthUsecase struct {
	appName    string
	oauthState time.Duration
	su         SessionUsecase
	ar         repositories.AuthRepository
	oar        repositories.OAuthRepository
	transactor *database.Transactor
	acp        *publisher.Publisher
	google     *oauth.OAuth
	logger     *logger.Logger
}

func NewOAuthUsecase(
	appName string,
	oauthState time.Duration,
	su SessionUsecase,
	ar repositories.AuthRepository,
	oar repositories.OAuthRepository,
	tx *database.Transactor,
	acp *publisher.Publisher,
	goa *oauth.OAuth,
	l *logger.Logger,
) OAuthUsecase {
	return &oAuthUsecase{
		appName:    appName,
		oauthState: oauthState,
		su:         su,
		ar:         ar,
		oar:        oar,
		transactor: tx,
		acp:        acp,
		google:     goa,
		logger:     l,
	}
}

func (u *oAuthUsecase) GetAuthURL(ctx context.Context, provider string) (string, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "oauth.usecase.GetAuthURL")
	defer span.End()

	providerField := logger.NewField("oauth_provider", provider)

	// OAuth Provider Selection
	oauth, err := u.getOAuthProvider(provider)
	if err != nil {
		return "", err.Append(providerField)
	}

	// CSRF State Creation
	// NOTE: Recommended length for CSRF state is 16 bytes
	state, genErr := utils.GenerateCSRFState(16)
	if genErr != nil {
		return "", ce.NewError(
			ce.CodeCSRFStateGenerationFailed,
			ce.MsgInternalServer,
			genErr,
			providerField,
		)
	}

	// CSRF State Storing
	err = u.oar.StoreState(
		ctx,
		&models.StoreOAuthState{
			State:    state,
			Provider: provider,
			Duration: u.oauthState,
		},
	)
	if err != nil {
		return "", err.Append(providerField)
	}

	// OAuth Consent Page URL Generation
	return oauth.AuthURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	), nil
}

func (u *oAuthUsecase) Callback(ctx context.Context, req *models.OAuthCallbackReq) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "oauth.usecase.Callback")
	defer span.End()

	providerField := logger.NewField("oauth_provider", req.Provider)

	// Data Normalization
	code := strings.TrimSpace(req.Code)
	state := strings.TrimSpace(req.State)

	// Data Validation
	if code == "" {
		return nil, ce.NewError(
			ce.CodeInvalidParams,
			"Authorization code is required",
			nil,
			providerField,
		)
	}
	if state == "" {
		return nil, ce.NewError(
			ce.CodeInvalidParams,
			"Authorization state is required",
			nil,
			providerField,
		)
	}

	// OAuth Provider Selection
	oauth, err := u.getOAuthProvider(req.Provider)
	if err != nil {
		return nil, err.Append(providerField)
	}

	// CSRF State Consumption
	provider, err := u.oar.UseState(ctx, state)
	if err != nil {
		return nil, err.Append(providerField)
	}
	if req.Provider != provider {
		return nil, ce.NewError(
			ce.CodeOAuthProviderMismatch,
			ce.MsgInvalidOAuthState,
			nil,
			providerField,
			logger.NewField("stored_provider", provider),
		)
	}

	// OAuth Code Exchange
	claim, exErr := oauth.Exchange(ctx, code)
	if exErr != nil {
		return nil, ce.NewError(
			ce.CodeOAuthCodeExchangeFailed,
			ce.MsgInternalServer,
			exErr,
			providerField,
		)
	}
	if !claim.IsEmailVerified {
		return nil, ce.NewError(
			ce.CodeEmailNotVerified,
			ce.MsgEmailNotVerified,
			nil,
			providerField,
		)
	}

	email := strings.ToLower(strings.TrimSpace(claim.Email))
	subject := strings.TrimSpace(claim.Subject)

	var a *models.Auth
	var at *models.AuthToken
	var isNewAccount bool
	err = u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Auth Fetching
		auth, err := u.ar.GetByEmail(ctx, email)
		if err != nil && err.Code() != ce.CodeAuthNotFound {
			return err
		}
		if err != nil {
			// Auth Creation
			now := time.Now().UTC()
			auth, err = u.ar.Create(
				ctx,
				&models.CreateAuth{
					Email:           email,
					Role:            constants.RoleUser,
					EmailVerifiedAt: &now,
				},
			)
			if err != nil {
				return err
			}

			// OAuth Creation
			_, err = u.oar.Create(
				ctx,
				&models.CreateOAuth{
					AuthID:   auth.ID,
					Provider: req.Provider,
					Subject:  subject,
				},
			)
			if err != nil {
				return err
			}

			isNewAccount = true
		}

		// Session Creation
		at, err = u.su.CreateSession(
			ctx,
			&models.CreateSessionReq{
				AuthID:          auth.ID,
				Role:            auth.Role,
				IsEmailVerified: auth.IsEmailVerified(),
			},
		)
		if err != nil {
			return err
		}

		a = auth
		return nil
	})
	if err != nil {
		return nil, err.Append(providerField)
	}

	authIDField := logger.NewField("auth_id", a.ID)

	// New Account Check
	if isNewAccount {
		evtTopicField := logger.NewField("event_topic", u.acp.Topic())

		// auth.created Event Publishing
		// NOTE: Fail to publish event does not fail Callback usecase
		err := u.acp.Publish(
			ctx,
			"auth_"+strconv.FormatUint(a.ID, 10),
			&events.AuthCreated{
				Id:              utils.GenerateUUID().String(),
				AuthId:          a.ID,
				Email:           a.Email,
				Role:            a.Role,
				EmailVerifiedAt: utils.ToTimestamp(a.EmailVerifiedAt),
				Session:         &at.RefreshToken.Token,
				CreatedAt:       timestamppb.New(time.Now().UTC()),
			},
		)
		if err != nil {
			u.logger.Warn(
				ctx,
				"created auth and oauth. failed to publish event",
				providerField,
				authIDField,
				evtTopicField,
				logger.NewField("error_code", ce.CodeEventPublishingFailed),
				logger.NewField("error", err),
			)
			return at, nil
		}

		u.logger.Info(
			ctx,
			"EVENT PUBLISHED",
			providerField,
			authIDField,
			evtTopicField,
		)
	}

	return at, nil
}

func (u *oAuthUsecase) getOAuthProvider(provider string) (*oauth.OAuth, *ce.Error) {
	switch provider {
	case constants.OAuthProviderGoogle:
		return u.google, nil
	default:
		return nil, ce.NewError(ce.CodeInvalidOAuthProvider, ce.MsgInvalidOAuthProvider, nil)
	}
}
