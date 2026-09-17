package repositories

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	"github.com/ritchieridanko/apotekly/services/product/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type ProductRepository interface {
	Create(ctx context.Context, data *models.CreateProduct) (p *models.Product, err *ce.Error)
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
