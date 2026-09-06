package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
	"github.com/ritchieridanko/apotekly/services/user/internal/repositories/database"
)

type AddressRepository interface {
	Create(ctx context.Context, data *models.CreateAddress) (a *models.Address, err *ce.Error)
	GetAll(ctx context.Context, authID uint64, params *models.GetAllAddresses) (as []models.Address, total int64, err *ce.Error)
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

func (r *addressRepository) GetAll(ctx context.Context, authID uint64, params *models.GetAllAddresses) ([]models.Address, int64, *ce.Error) {
	return r.database.GetAll(ctx, authID, params)
}

func (r *addressRepository) UnsetPrimary(ctx context.Context, authID uint64) (*models.Address, *ce.Error) {
	return r.database.UnsetPrimary(ctx, authID)
}
