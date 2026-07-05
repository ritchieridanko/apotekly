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
	GetByEmail(ctx context.Context, email string) (a *models.Auth, err *ce.Error)
	GetByID(ctx context.Context, authID uint64) (a *models.Auth, err *ce.Error)
	UpdatePassword(ctx context.Context, authID uint64, newPassword string) (err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
	SetVerified(ctx context.Context, authID uint64) (a *models.Auth, err *ce.Error)
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

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*models.Auth, *ce.Error) {
	return r.database.GetByEmail(ctx, email)
}

func (r *authRepository) GetByID(ctx context.Context, authID uint64) (*models.Auth, *ce.Error) {
	return r.database.GetByID(ctx, authID)
}

func (r *authRepository) UpdatePassword(ctx context.Context, authID uint64, newPassword string) *ce.Error {
	return r.database.UpdatePassword(ctx, authID, newPassword)
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

func (r *authRepository) SetVerified(ctx context.Context, authID uint64) (*models.Auth, *ce.Error) {
	return r.database.SetVerified(ctx, authID)
}
