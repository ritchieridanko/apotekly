package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/models"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type PharmacyRepository interface {
	Create(ctx context.Context, data *models.CreatePharmacy) (p *models.Pharmacy, err *ce.Error)
	GetID(ctx context.Context, authID uint64) (pharmacyID uuid.UUID, err *ce.Error)
	GetByAuthID(ctx context.Context, authID uint64) (p *models.Pharmacy, err *ce.Error)
	Update(ctx context.Context, authID uint64, data *models.UpdatePharmacy) (p *models.Pharmacy, err *ce.Error)
}

type pharmacyRepository struct {
	database database.PharmacyDatabase
}

func NewPharmacyRepository(db database.PharmacyDatabase) PharmacyRepository {
	return &pharmacyRepository{database: db}
}

func (r *pharmacyRepository) Create(ctx context.Context, data *models.CreatePharmacy) (*models.Pharmacy, *ce.Error) {
	return r.database.Create(ctx, data)
}

func (r *pharmacyRepository) GetID(ctx context.Context, authID uint64) (uuid.UUID, *ce.Error) {
	return r.database.GetID(ctx, authID)
}

func (r *pharmacyRepository) GetByAuthID(ctx context.Context, authID uint64) (*models.Pharmacy, *ce.Error) {
	return r.database.GetByAuthID(ctx, authID)
}

func (r *pharmacyRepository) Update(ctx context.Context, authID uint64, data *models.UpdatePharmacy) (*models.Pharmacy, *ce.Error) {
	return r.database.Update(ctx, authID, data)
}
