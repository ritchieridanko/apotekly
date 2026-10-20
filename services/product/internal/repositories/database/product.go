package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type ProductDatabase interface {
	Create(ctx context.Context, data *models.CreateProduct) (p *models.Product, err *ce.Error)
	GetPharmacyID(ctx context.Context, productID uuid.UUID) (pharmacyID uuid.UUID, err *ce.Error)
	GetCurrency(ctx context.Context, productID uuid.UUID) (currency string, err *ce.Error)
	GetByID(ctx context.Context, productID uuid.UUID) (p *models.Product, err *ce.Error)
	Update(ctx context.Context, params *models.UpdateProductP, data *models.UpdateProductD) (p *models.Product, err *ce.Error)
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

func (d *productDatabase) GetPharmacyID(ctx context.Context, productID uuid.UUID) (uuid.UUID, *ce.Error) {
	query := "SELECT pharmacy_id FROM products WHERE id = $1 AND deleted_at IS NULL"
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var pharmacyID uuid.UUID
	err := d.database.Query(
		ctx, query,
		productID,
	).Scan(
		&pharmacyID,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get product pharmacy id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return uuid.Nil, ce.NewError(
				ce.CodeProductNotFound,
				ce.MsgProductNotFound,
				wrappedErr,
			)
		}
		return uuid.Nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return pharmacyID, nil
}

func (d *productDatabase) GetCurrency(ctx context.Context, productID uuid.UUID) (string, *ce.Error) {
	query := "SELECT currency FROM products WHERE id = $1 AND deleted_at IS NULL"
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var currency string
	err := d.database.Query(
		ctx, query,
		productID,
	).Scan(
		&currency,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get product currency: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return "", ce.NewError(
				ce.CodeProductNotFound,
				ce.MsgProductNotFound,
				wrappedErr,
			)
		}
		return "", ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return currency, nil
}

func (d *productDatabase) GetByID(ctx context.Context, productID uuid.UUID) (*models.Product, *ce.Error) {
	query := `
		SELECT
			id, pharmacy_id, brand_name, generic_name, description, requires_rx,
			dosage_form, strength, pack_unit, pack_size, height_cm, length_cm,
			width_cm, weight_g, price, currency, quantity, is_active, manufactured_by,
			manufactured_in, reg_authority, reg_identifier, created_at, updated_at
		FROM
			products
		WHERE
			id = $1
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var p models.Product
	err := d.database.Query(
		ctx, query,
		productID,
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
		wrappedErr := fmt.Errorf("failed to get product by id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeProductNotFound,
				ce.MsgProductNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &p, nil
}

func (d *productDatabase) Update(ctx context.Context, params *models.UpdateProductP, data *models.UpdateProductD) (*models.Product, *ce.Error) {
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if data.BrandName != nil {
		setClauses = append(setClauses, "brand_name = $"+strconv.Itoa(argPos))
		args = append(args, *data.BrandName)
		argPos++
	}
	if data.GenericName != nil {
		setClauses = append(setClauses, "generic_name = $"+strconv.Itoa(argPos))
		args = append(args, *data.GenericName)
		argPos++
	}
	if data.Description != nil {
		setClauses = append(setClauses, "description = $"+strconv.Itoa(argPos))
		args = append(args, *data.Description)
		argPos++
	}
	if data.RequiresRX != nil {
		setClauses = append(setClauses, "requires_rx = $"+strconv.Itoa(argPos))
		args = append(args, *data.RequiresRX)
		argPos++
	}
	if data.DosageForm != nil {
		setClauses = append(setClauses, "dosage_form = $"+strconv.Itoa(argPos))
		args = append(args, *data.DosageForm)
		argPos++
	}
	if data.Strength != nil {
		setClauses = append(setClauses, "strength = $"+strconv.Itoa(argPos))
		args = append(args, *data.Strength)
		argPos++
	}
	if data.PackUnit != nil {
		setClauses = append(setClauses, "pack_unit = $"+strconv.Itoa(argPos))
		args = append(args, *data.PackUnit)
		argPos++
	}
	if data.PackSize != nil {
		setClauses = append(setClauses, "pack_size = $"+strconv.Itoa(argPos))
		args = append(args, *data.PackSize)
		argPos++
	}
	if data.HeightCM != nil {
		setClauses = append(setClauses, "height_cm = $"+strconv.Itoa(argPos))
		args = append(args, *data.HeightCM)
		argPos++
	}
	if data.LengthCM != nil {
		setClauses = append(setClauses, "length_cm = $"+strconv.Itoa(argPos))
		args = append(args, *data.LengthCM)
		argPos++
	}
	if data.WidthCM != nil {
		setClauses = append(setClauses, "width_cm = $"+strconv.Itoa(argPos))
		args = append(args, *data.WidthCM)
		argPos++
	}
	if data.WeightG != nil {
		setClauses = append(setClauses, "weight_g = $"+strconv.Itoa(argPos))
		args = append(args, *data.WeightG)
		argPos++
	}
	if data.Price != nil {
		setClauses = append(setClauses, "price = $"+strconv.Itoa(argPos))
		args = append(args, *data.Price)
		argPos++
	}
	if data.Currency != nil {
		setClauses = append(setClauses, "currency = $"+strconv.Itoa(argPos))
		args = append(args, *data.Currency)
		argPos++
	}
	if data.Quantity != nil {
		setClauses = append(setClauses, "quantity = $"+strconv.Itoa(argPos))
		args = append(args, *data.Quantity)
		argPos++
	}
	if data.ManufacturedBy != nil {
		setClauses = append(setClauses, "manufactured_by = $"+strconv.Itoa(argPos))
		args = append(args, *data.ManufacturedBy)
		argPos++
	}
	if data.ManufacturedIn != nil {
		setClauses = append(setClauses, "manufactured_in = $"+strconv.Itoa(argPos))
		args = append(args, *data.ManufacturedIn)
		argPos++
	}
	if data.RegAuthority != nil {
		setClauses = append(setClauses, "reg_authority = $"+strconv.Itoa(argPos))
		args = append(args, *data.RegAuthority)
		argPos++
	}
	if data.RegIdentifier != nil {
		setClauses = append(setClauses, "reg_identifier = $"+strconv.Itoa(argPos))
		args = append(args, *data.RegIdentifier)
		argPos++
	}
	if len(setClauses) == 0 {
		return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, nil)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, params.ProductID, params.PharmacyID)
	query := fmt.Sprintf(
		`
			UPDATE
				products
			SET
				%s
			WHERE
				id = $%d
				AND pharmacy_id = $%d
				AND deleted_at IS NULL
			RETURNING
				id, pharmacy_id, brand_name, generic_name, description, requires_rx,
				dosage_form, strength, pack_unit, pack_size, height_cm, length_cm,
				width_cm, weight_g, price, currency, quantity, is_active, manufactured_by,
				manufactured_in, reg_authority, reg_identifier, created_at, updated_at
		`,
		strings.Join(setClauses, ", "), argPos, argPos+1,
	)

	var p models.Product
	err := d.database.Query(
		ctx, query,
		args...,
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
		wrappedErr := fmt.Errorf("failed to update product: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeProductNotFound,
				ce.MsgProductNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &p, nil
}
