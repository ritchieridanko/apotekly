package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
	"github.com/ritchieridanko/apotekly/services/user/internal/repositories/database"
)

type UserRepository interface {
	Create(ctx context.Context, data *models.CreateUser) (u *models.User, err *ce.Error)
	GetID(ctx context.Context, authID uint64) (userID uuid.UUID, err *ce.Error)
	GetByAuthID(ctx context.Context, authID uint64) (u *models.User, err *ce.Error)
	Update(ctx context.Context, authID uint64, data *models.UpdateUser) (u *models.User, err *ce.Error)
}

type userRepository struct {
	database database.UserDatabase
}

func NewUserRepository(db database.UserDatabase) UserRepository {
	return &userRepository{database: db}
}

func (r *userRepository) Create(ctx context.Context, data *models.CreateUser) (*models.User, *ce.Error) {
	return r.database.Create(ctx, data)
}

func (r *userRepository) GetID(ctx context.Context, authID uint64) (uuid.UUID, *ce.Error) {
	return r.database.GetID(ctx, authID)
}

func (r *userRepository) GetByAuthID(ctx context.Context, authID uint64) (*models.User, *ce.Error) {
	return r.database.GetByAuthID(ctx, authID)
}

func (r *userRepository) Update(ctx context.Context, authID uint64, data *models.UpdateUser) (*models.User, *ce.Error) {
	return r.database.Update(ctx, authID, data)
}
