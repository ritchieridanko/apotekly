package usecases

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/contract/events/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/publisher"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/bcrypt"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, req *models.SignUpReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	SignIn(ctx context.Context, req *models.SignInReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
}

type authUsecase struct {
	appName           string
	verificationToken time.Duration
	su                SessionUsecase
	ar                repositories.AuthRepository
	tr                repositories.TokenRepository
	transactor        *database.Transactor
	acp               *publisher.Publisher
	bcrypt            *bcrypt.BCrypt
	validator         *validator.Validator
	logger            *logger.Logger
}

func NewAuthUsecase(
	appName string,
	verificationToken time.Duration,
	su SessionUsecase,
	ar repositories.AuthRepository,
	tr repositories.TokenRepository,
	tx *database.Transactor,
	acp *publisher.Publisher,
	b *bcrypt.BCrypt,
	v *validator.Validator,
	l *logger.Logger,
) AuthUsecase {
	return &authUsecase{
		appName:           appName,
		verificationToken: verificationToken,
		su:                su,
		ar:                ar,
		tr:                tr,
		transactor:        tx,
		acp:               acp,
		bcrypt:            b,
		validator:         v,
		logger:            l,
	}
}

func (u *authUsecase) SignUp(ctx context.Context, req *models.SignUpReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.SignUp")
	defer span.End()

	// Data Normalization
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Data Validation
	if ok, why := u.validator.Email(email); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}
	if ok, why := u.validator.Password(req.Password); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	var a *models.Auth
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Email Availability Check
		available, err := u.ar.IsEmailAvailable(ctx, email)
		if err != nil {
			return err
		}
		if !available {
			return ce.NewError(ce.CodeEmailNotAvailable, ce.MsgEmailAlreadyRegistered, nil)
		}

		// Password Hashing
		hash, hashErr := u.bcrypt.Hash(req.Password)
		if hashErr != nil {
			return ce.NewError(ce.CodeBCryptHashingFailed, ce.MsgInternalServer, hashErr)
		}

		// Auth Creation
		a, err = u.ar.Create(
			ctx,
			&models.CreateAuth{
				Email:    email,
				Password: &hash,
				Role:     constants.RoleUser,
			},
		)
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	authIDField := logger.NewField("auth_id", a.ID)

	// Session Creation
	// NOTE: Fail to create session does not fail SignUp usecase
	var session *string
	at, err := u.su.CreateSession(
		ctx,
		&models.CreateSessionReq{
			AuthID:          a.ID,
			Role:            a.Role,
			IsEmailVerified: a.IsEmailVerified(),
		},
	)
	if err != nil {
		u.logger.Warn(
			ctx,
			"created auth. failed to create session",
			err.Append(
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err.Unwrap()),
			).Fields()...,
		)
	} else {
		session = &at.RefreshToken.Token
	}

	// Verification Token Creation
	// NOTE: Fail to create verification token does not fail SignUp usecase
	token := utils.GenerateUUID().String()
	verificationToken := &token
	err = u.tr.CreateVerification(
		ctx,
		&models.CreateVerificationToken{
			AuthID:   a.ID,
			Token:    token,
			Duration: u.verificationToken,
		},
	)
	if err != nil {
		verificationToken = nil
		u.logger.Warn(
			ctx,
			"created auth. failed to create verification token",
			err.Append(
				authIDField,
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err.Unwrap()),
			).Fields()...,
		)
	}

	evtTopicField := logger.NewField("event_topic", u.acp.Topic())

	// auth.created Event Publishing
	// NOTE: Fail to publish event does not fail SignUp usecase
	pubErr := u.acp.Publish(
		ctx,
		"auth_"+strconv.FormatUint(a.ID, 10),
		&events.AuthCreated{
			Id:                utils.GenerateUUID().String(),
			AuthId:            a.ID,
			Email:             a.Email,
			Role:              a.Role,
			EmailVerifiedAt:   utils.ToTimestamp(a.EmailVerifiedAt),
			Session:           session,
			VerificationToken: verificationToken,
			CreatedAt:         timestamppb.New(time.Now().UTC()),
		},
	)
	if pubErr != nil {
		u.logger.Warn(
			ctx,
			"created auth. failed to publish event",
			authIDField,
			evtTopicField,
			logger.NewField("error_code", ce.CodeEventPublishingFailed),
			logger.NewField("error", pubErr),
		)
		return a, at, nil
	}

	u.logger.Info(
		ctx,
		"EVENT PUBLISHED",
		authIDField,
		evtTopicField,
	)

	return a, at, nil
}

func (u *authUsecase) SignIn(ctx context.Context, req *models.SignInReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.SignIn")
	defer span.End()

	// Data Normalization
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Data Validation
	if ok, why := u.validator.Email(email); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	// Auth Fetching
	a, err := u.ar.GetByEmail(ctx, email)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return nil, nil, ce.NewError(
			ce.CodeEmailNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
			err.Fields()...,
		)
	}
	if err != nil {
		return nil, nil, err
	}

	authIDField := logger.NewField("auth_id", a.ID)

	// Password Validation
	if a.Password == nil {
		return nil, nil, ce.NewError(
			ce.CodeOAuthRegularSignIn,
			ce.MsgInvalidCredentials,
			errors.New("auth has no password"),
			authIDField,
		)
	}
	if err := u.bcrypt.Validate(*a.Password, req.Password); err != nil {
		return nil, nil, ce.NewError(
			ce.CodeWrongPassword,
			ce.MsgInvalidCredentials,
			err,
			authIDField,
		)
	}

	// Session Creation
	at, err := u.su.CreateSession(
		ctx,
		&models.CreateSessionReq{
			AuthID:          a.ID,
			Role:            a.Role,
			IsEmailVerified: a.IsEmailVerified(),
		},
	)

	return a, at, err
}

func (u *authUsecase) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.IsEmailAvailable")
	defer span.End()

	// Data Normalization
	em := strings.ToLower(strings.TrimSpace(email))

	// Data Validation
	if ok, why := u.validator.Email(em); !ok {
		return false, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	// Email Availability Check
	return u.ar.IsEmailAvailable(ctx, em)
}
