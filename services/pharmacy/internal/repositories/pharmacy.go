package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/models"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type PharmacyRepository interface {
	Create(ctx context.Context, data *models.CreatePharmacy) (p *models.Pharmacy, err *ce.Error)
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
