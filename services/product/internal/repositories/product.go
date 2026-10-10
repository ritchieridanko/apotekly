package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	"github.com/ritchieridanko/apotekly/services/product/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type ProductRepository interface {
	Create(ctx context.Context, data *models.CreateProduct) (p *models.Product, err *ce.Error)
	GetPharmacyID(ctx context.Context, productID uuid.UUID) (pharmacyID uuid.UUID, err *ce.Error)
}

type productRepository struct {
	database database.ProductDatabase
}

func NewProductRepository(db database.ProductDatabase) ProductRepository {
	return &productRepository{database: db}
}

func (r *productRepository) Create(ctx context.Context, data *models.CreateProduct) (*models.Product, *ce.Error) {
	return r.database.Create(ctx, data)
}

func (r *productRepository) GetPharmacyID(ctx context.Context, productID uuid.UUID) (uuid.UUID, *ce.Error) {
	return r.database.GetPharmacyID(ctx, productID)
}
