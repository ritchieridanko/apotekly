package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/cache"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type AuthRepository interface {
	Create(ctx context.Context, data *models.CreateAuth) (a *models.Auth, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
}

type authRepository struct {
	cache    cache.AuthCache
	database database.AuthDatabase
}

func NewAuthRepository(cc cache.AuthCache, db database.AuthDatabase) AuthRepository {
	return &authRepository{cache: cc, database: db}
}

func (r *authRepository) Create(ctx context.Context, data *models.CreateAuth) (*models.Auth, *ce.Error) {
	return r.database.Create(ctx, data)
}

func (r *authRepository) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	exists, err := r.database.EmailExists(ctx, email)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	reserved, err := r.cache.IsEmailReserved(ctx, email)
	if err != nil {
		return false, err
	}
	return !reserved, nil
}
