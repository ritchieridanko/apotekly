package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	"github.com/ritchieridanko/apotekly/services/product/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type ProductImageRepository interface {
	CreateMany(ctx context.Context, data *models.CreateProductImages) (pis []models.ProductImage, err *ce.Error)
	GetAll(ctx context.Context, productID uuid.UUID) (pis []models.ProductImage, err *ce.Error)
	Update(ctx context.Context, params *models.UpdateProductImageP, data *models.UpdateProductImageD) (pi *models.ProductImage, err *ce.Error)
	CountAll(ctx context.Context, productID uuid.UUID) (total int64, highestSortOrder int16, err *ce.Error)

	// NOTE: Shall be placed WITHIN a transaction
	Delete(ctx context.Context, params *models.DeleteProductImage) (err *ce.Error)
	Reorder(ctx context.Context, params *models.ReorderProductImage, newSortOrder int16) (pi *models.ProductImage, err *ce.Error)
}

type productImageRepository struct {
	database database.ProductImageDatabase
}

func NewProductImageRepository(db database.ProductImageDatabase) ProductImageRepository {
	return &productImageRepository{database: db}
}

func (r *productImageRepository) CreateMany(ctx context.Context, data *models.CreateProductImages) ([]models.ProductImage, *ce.Error) {
	return r.database.CreateMany(ctx, data)
}

func (r *productImageRepository) GetAll(ctx context.Context, productID uuid.UUID) ([]models.ProductImage, *ce.Error) {
	return r.database.GetAll(ctx, productID)
}

func (r *productImageRepository) Update(ctx context.Context, params *models.UpdateProductImageP, data *models.UpdateProductImageD) (*models.ProductImage, *ce.Error) {
	return r.database.Update(ctx, params, data)
}

func (r *productImageRepository) Delete(ctx context.Context, params *models.DeleteProductImage) *ce.Error {
	return r.database.Delete(ctx, params)
}

func (r *productImageRepository) CountAll(ctx context.Context, productID uuid.UUID) (int64, int16, *ce.Error) {
	return r.database.CountAll(ctx, productID)
}

func (r *productImageRepository) Reorder(ctx context.Context, params *models.ReorderProductImage, newSortOrder int16) (*models.ProductImage, *ce.Error) {
	return r.database.Reorder(ctx, params, newSortOrder)
}
