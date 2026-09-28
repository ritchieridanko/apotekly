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
	GetActiveStatus(ctx context.Context, authID uint64) (active bool, err *ce.Error)
	GetByID(ctx context.Context, pharmacyID uuid.UUID) (p *models.Pharmacy, err *ce.Error)
	GetByAuthID(ctx context.Context, authID uint64) (p *models.Pharmacy, err *ce.Error)
	GetAll(ctx context.Context, params *models.GetAllPharmacies) (pss []models.PharmacySummary, total int64, err *ce.Error)
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

func (r *pharmacyRepository) GetActiveStatus(ctx context.Context, authID uint64) (bool, *ce.Error) {
	return r.database.GetActiveStatus(ctx, authID)
}

func (r *pharmacyRepository) GetByID(ctx context.Context, pharmacyID uuid.UUID) (*models.Pharmacy, *ce.Error) {
	return r.database.GetByID(ctx, pharmacyID)
}

func (r *pharmacyRepository) GetByAuthID(ctx context.Context, authID uint64) (*models.Pharmacy, *ce.Error) {
	return r.database.GetByAuthID(ctx, authID)
}

func (r *pharmacyRepository) GetAll(ctx context.Context, params *models.GetAllPharmacies) ([]models.PharmacySummary, int64, *ce.Error) {
	return r.database.GetAll(ctx, params)
}

func (r *pharmacyRepository) Update(ctx context.Context, authID uint64, data *models.UpdatePharmacy) (*models.Pharmacy, *ce.Error) {
	return r.database.Update(ctx, authID, data)
}
