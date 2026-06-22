package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/cache"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type TokenRepository interface {
	CreateVerification(ctx context.Context, data *models.CreateVerificationToken) (err *ce.Error)
}

type tokenRepository struct {
	cache cache.TokenCache
}

func NewTokenRepository(cc cache.TokenCache) TokenRepository {
	return &tokenRepository{cache: cc}
}

func (r *tokenRepository) CreateVerification(ctx context.Context, data *models.CreateVerificationToken) *ce.Error {
	return r.cache.CreateVerification(ctx, data)
}
