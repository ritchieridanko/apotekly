package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/cache"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type TokenRepository interface {
	CreateEmailChange(ctx context.Context, data *models.CreateEmailChangeToken) (err *ce.Error)
	UseEmailChange(ctx context.Context, token string) (authID uint64, newEmail string, err *ce.Error)
	CreateVerification(ctx context.Context, data *models.CreateVerificationToken) (err *ce.Error)
	UseVerification(ctx context.Context, token string) (authID uint64, err *ce.Error)
}

type tokenRepository struct {
	cache cache.TokenCache
}

func NewTokenRepository(cc cache.TokenCache) TokenRepository {
	return &tokenRepository{cache: cc}
}

func (r *tokenRepository) CreateEmailChange(ctx context.Context, data *models.CreateEmailChangeToken) *ce.Error {
	return r.cache.CreateEmailChange(ctx, data)
}

func (r *tokenRepository) UseEmailChange(ctx context.Context, token string) (uint64, string, *ce.Error) {
	return r.cache.UseEmailChange(ctx, token)
}

func (r *tokenRepository) CreateVerification(ctx context.Context, data *models.CreateVerificationToken) *ce.Error {
	return r.cache.CreateVerification(ctx, data)
}

func (r *tokenRepository) UseVerification(ctx context.Context, token string) (uint64, *ce.Error) {
	return r.cache.UseVerification(ctx, token)
}
