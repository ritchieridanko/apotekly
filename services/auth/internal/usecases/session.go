package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/jwt"
	"go.opentelemetry.io/otel"
)

type SessionUsecase interface {
	CreateSession(ctx context.Context, req *models.CreateSessionReq) (at *models.AuthToken, err *ce.Error)
}

type sessionUsecase struct {
	appName      string
	accessToken  time.Duration
	refreshToken time.Duration
	sr           repositories.SessionRepository
	transactor   *database.Transactor
	jwt          *jwt.JWT
}

func NewSessionUsecase(
	appName string,
	accessToken,
	refreshToken time.Duration,
	sr repositories.SessionRepository,
	tx *database.Transactor,
	j *jwt.JWT,
) SessionUsecase {
	return &sessionUsecase{
		appName:      appName,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		sr:           sr,
		transactor:   tx,
		jwt:          j,
	}
}

func (u *sessionUsecase) CreateSession(ctx context.Context, req *models.CreateSessionReq) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "session.usecase.CreateSession")
	defer span.End()

	authIDField := logger.NewField("auth_id", req.AuthID)

	// UUID Creation
	uuid := utils.GenerateUUID()

	// JWT Creation
	now := time.Now().UTC()
	jwt, err := u.jwt.Generate(req.AuthID, req.Role, req.IsEmailVerified, &now)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeJWTGenerationFailed,
			ce.MsgInternalServer,
			err,
			authIDField,
		)
	}

	txErr := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		transportCtx := utils.CtxTransport(ctx)
		if transportCtx == nil {
			return ce.NewError(
				ce.CodeMissingContextValue,
				ce.MsgInternalServer,
				errors.New("transport missing from context"),
			)
		}

		// Active Session Revocation
		sessionID, err := u.sr.RevokeActive(
			ctx,
			&models.RevokeActiveSession{
				AuthID:    req.AuthID,
				IPAddress: transportCtx.IPAddress,
				UserAgent: transportCtx.UserAgent,
				ExpiresAt: now,
			},
		)
		if err != nil {
			return err
		}

		// Session Creation
		// NOTE: Set revoked session (if any) as parent session to maintain session lineage
		data := models.CreateSession{
			AuthID:       req.AuthID,
			RefreshToken: uuid.String(),
			IPAddress:    transportCtx.IPAddress,
			UserAgent:    transportCtx.UserAgent,
			ExpiresAt:    now.Add(u.refreshToken),
		}
		if invalidSessionID := uint64(0); sessionID != invalidSessionID {
			data.ParentID = &sessionID
		}
		return u.sr.Create(ctx, &data)
	})
	if txErr != nil {
		return nil, txErr.Append(authIDField)
	}

	return &models.AuthToken{
		AccessToken: &models.AccessToken{
			Token:            jwt,
			ExpiresInSeconds: uint64(u.accessToken.Seconds()),
		},
		RefreshToken: &models.RefreshToken{
			Token:            uuid.String(),
			ExpiresInSeconds: uint64(u.refreshToken.Seconds()),
		},
	}, nil
}
