package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type SessionRepository interface {
	Create(ctx context.Context, data *models.CreateSession) (err *ce.Error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (s *models.Session, err *ce.Error)
	Revoke(ctx context.Context, params *models.RevokeSession) (s *models.Session, err *ce.Error)
	RevokeActive(ctx context.Context, params *models.RevokeActiveSession) (sessionID uint64, err *ce.Error)
	RevokeAllActive(ctx context.Context, params *models.RevokeAllActiveSessions) (err *ce.Error)
}

type sessionRepository struct {
	database database.SessionDatabase
}

func NewSessionRepository(db database.SessionDatabase) SessionRepository {
	return &sessionRepository{database: db}
}

func (r *sessionRepository) Create(ctx context.Context, data *models.CreateSession) *ce.Error {
	return r.database.Create(ctx, data)
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, *ce.Error) {
	return r.database.GetByRefreshToken(ctx, refreshToken)
}

func (r *sessionRepository) Revoke(ctx context.Context, params *models.RevokeSession) (*models.Session, *ce.Error) {
	return r.database.Revoke(ctx, params)
}

func (r *sessionRepository) RevokeActive(ctx context.Context, params *models.RevokeActiveSession) (uint64, *ce.Error) {
	return r.database.RevokeActive(ctx, params)
}

func (r *sessionRepository) RevokeAllActive(ctx context.Context, params *models.RevokeAllActiveSessions) *ce.Error {
	return r.database.RevokeAllActive(ctx, params)
}
