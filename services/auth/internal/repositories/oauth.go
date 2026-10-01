package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/cache"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type OAuthRepository interface {
	Create(ctx context.Context, data *models.CreateOAuth) (oa *models.OAuth, err *ce.Error)
	StoreState(ctx context.Context, data *models.StoreOAuthState) (err *ce.Error)
	UseState(ctx context.Context, state string) (provider string, err *ce.Error)
}

type oAuthRepository struct {
	cache    cache.OAuthCache
	database database.OAuthDatabase
}

func NewOAuthRepository(cc cache.OAuthCache, db database.OAuthDatabase) OAuthRepository {
	return &oAuthRepository{cache: cc, database: db}
}

func (r *oAuthRepository) Create(ctx context.Context, data *models.CreateOAuth) (*models.OAuth, *ce.Error) {
	return r.database.Create(ctx, data)
}

func (r *oAuthRepository) StoreState(ctx context.Context, data *models.StoreOAuthState) *ce.Error {
	return r.cache.StoreState(ctx, data)
}

func (r *oAuthRepository) UseState(ctx context.Context, state string) (string, *ce.Error) {
	return r.cache.UseState(ctx, state)
}
