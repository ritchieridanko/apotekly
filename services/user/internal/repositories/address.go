package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
	"github.com/ritchieridanko/apotekly/services/user/internal/repositories/database"
)

type AddressRepository interface {
	Create(ctx context.Context, data *models.CreateAddress) (a *models.Address, err *ce.Error)
	UnsetPrimary(ctx context.Context, authID uint64) (a *models.Address, err *ce.Error)
}

type addressRepository struct {
	database database.AddressDatabase
}

func NewAddressRepository(db database.AddressDatabase) AddressRepository {
	return &addressRepository{database: db}
}

func (r *addressRepository) Create(ctx context.Context, data *models.CreateAddress) (*models.Address, *ce.Error) {
	return r.database.Create(ctx, data)
}

func (r *addressRepository) UnsetPrimary(ctx context.Context, authID uint64) (*models.Address, *ce.Error) {
	return r.database.UnsetPrimary(ctx, authID)
}
