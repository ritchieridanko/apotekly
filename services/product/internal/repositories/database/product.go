package database

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type ProductDatabase interface {
	Create(ctx context.Context, data *models.CreateProduct) (p *models.Product, err *ce.Error)
}

type productDatabase struct {
	database *db.Database
}

func NewProductDatabase(db *db.Database) ProductDatabase {
	return &productDatabase{database: db}
}

func (d *productDatabase) Create(ctx context.Context, data *models.CreateProduct) (*models.Product, *ce.Error) {
	query := `
		INSERT INTO products (
			id, pharmacy_id, brand_name, generic_name, description, requires_rx,
			dosage_form, strength, pack_unit, pack_size, height_cm, length_cm,
			width_cm, weight_g, price, currency, quantity, is_active,
			manufactured_by, manufactured_in, reg_authority, reg_identifier
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
		RETURNING
			id, pharmacy_id, brand_name, generic_name, description, requires_rx,
			dosage_form, strength, pack_unit, pack_size, height_cm, length_cm,
			width_cm, weight_g, price, currency, quantity, is_active, manufactured_by,
			manufactured_in, reg_authority, reg_identifier, created_at, updated_at
	`

	var p models.Product
	err := d.database.Query(
		ctx, query,
		data.ID,
		data.PharmacyID,
		data.BrandName,
		data.GenericName,
		data.Description,
		data.RequiresRX,
		data.DosageForm,
		data.Strength,
		data.PackUnit,
		data.PackSize,
		data.HeightCM,
		data.LengthCM,
		data.WidthCM,
		data.WeightG,
		data.Price,
		data.Currency,
		data.Quantity,
		data.IsActive,
		data.ManufacturedBy,
		data.ManufacturedIn,
		data.RegAuthority,
		data.RegIdentifier,
	).Scan(
		&p.ID,
		&p.PharmacyID,
		&p.BrandName,
		&p.GenericName,
		&p.Description,
		&p.RequiresRX,
		&p.DosageForm,
		&p.Strength,
		&p.PackUnit,
		&p.PackSize,
		&p.HeightCM,
		&p.LengthCM,
		&p.WidthCM,
		&p.WeightG,
		&p.Price,
		&p.Currency,
		&p.Quantity,
		&p.IsActive,
		&p.ManufacturedBy,
		&p.ManufacturedIn,
		&p.RegAuthority,
		&p.RegIdentifier,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create product: %w", err),
		)
	}

	return &p, nil
}
