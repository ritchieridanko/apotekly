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

type ProductImageDatabase interface {
	CreateMany(ctx context.Context, data *models.CreateProductImages) (pis []models.ProductImage, err *ce.Error)
	Update(ctx context.Context, params *models.UpdateProductImageP, data *models.UpdateProductImageD) (pi *models.ProductImage, err *ce.Error)
	CountAll(ctx context.Context, productID uuid.UUID) (total int64, highestSortOrder int16, err *ce.Error)
}

type productImageDatabase struct {
	database *db.Database
}

func NewProductImageDatabase(db *db.Database) ProductImageDatabase {
	return &productImageDatabase{database: db}
}

func (d *productImageDatabase) CreateMany(ctx context.Context, data *models.CreateProductImages) ([]models.ProductImage, *ce.Error) {
	totalImages := len(data.Images)
	if totalImages == 0 {
		return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, nil)
	}

	valClauses := []string{}
	args := []any{}
	numCols := 6

	for i, image := range data.Images {
		base := i * numCols

		valClauses = append(valClauses,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d)",
				base+1, base+2, base+3, base+4, base+5, base+6,
			),
		)
		args = append(args,
			image.ID,
			data.ProductID,
			image.Name,
			image.URL,
			image.AltText,
			image.SortOrder,
		)
	}

	query := fmt.Sprintf(
		`
			INSERT INTO product_images (
				id, product_id, name, image, alt_text, sort_order
			)
			VALUES
				%s
			RETURNING
				id, product_id, name, image, alt_text, sort_order
		`,
		strings.Join(valClauses, ", "),
	)

	rows, err := d.database.QueryAll(
		ctx, query,
		args...,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create product images: %w", err),
		)
	}
	defer rows.Close()

	pis := make([]models.ProductImage, 0, totalImages)

	for rows.Next() {
		var pi models.ProductImage
		err := rows.Scan(
			&pi.ID,
			&pi.ProductID,
			&pi.Name,
			&pi.Image,
			&pi.AltText,
			&pi.SortOrder,
		)
		if err != nil {
			return nil, ce.NewError(
				ce.CodeDBQueryExec,
				ce.MsgInternalServer,
				fmt.Errorf("failed to create product images: %w", err),
			)
		}

		pis = append(pis, pi)
	}
	if err := rows.Err(); err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create product images: %w", err),
		)
	}

	return pis, nil
}

func (d *productImageDatabase) Update(ctx context.Context, params *models.UpdateProductImageP, data *models.UpdateProductImageD) (*models.ProductImage, *ce.Error) {
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if data.Name != nil {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argPos))
		args = append(args, *data.Name)
		argPos++
	}
	if data.AltText != nil {
		setClauses = append(setClauses, "alt_text = $"+strconv.Itoa(argPos))
		args = append(args, *data.AltText)
		argPos++
	}
	if len(setClauses) == 0 {
		return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, nil)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, params.ProductImageID, params.ProductID)
	query := fmt.Sprintf(
		`
			UPDATE
				product_images
			SET
				%s
			WHERE
				id = $%d
				AND product_id = $%d
			RETURNING
				id, product_id, name, image, alt_text, sort_order
		`,
		strings.Join(setClauses, ", "), argPos, argPos+1,
	)

	var pi models.ProductImage
	err := d.database.Query(
		ctx, query,
		args...,
	).Scan(
		&pi.ID,
		&pi.ProductID,
		&pi.Name,
		&pi.Image,
		&pi.AltText,
		&pi.SortOrder,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product image: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeProductImageNotFound,
				ce.MsgProductImageNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &pi, nil
}

func (d *productImageDatabase) CountAll(ctx context.Context, productID uuid.UUID) (int64, int16, *ce.Error) {
	query := `
		SELECT
			sort_order,

			COUNT(*) OVER() AS total
		FROM
			product_images
		WHERE
			product_id = $1
		ORDER BY
			sort_order DESC
		LIMIT
			1
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var sortOrder int16
	var total int64
	err := d.database.Query(
		ctx, query,
		productID,
	).Scan(
		&sortOrder,
		&total,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, 0, nil
		}
		return 0, 0, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to count all product images: %w", err),
		)
	}

	return total, sortOrder, nil
}
