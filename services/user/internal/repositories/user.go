package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
	"github.com/ritchieridanko/apotekly/services/user/internal/repositories/database"
)

type UserRepository interface {
	Create(ctx context.Context, data *models.CreateUser) (u *models.User, err *ce.Error)
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
