package database

import (
	"context"
	"errors"
	"fmt"

	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
)

type AddressDatabase interface {
	Create(ctx context.Context, data *models.CreateAddress) (a *models.Address, err *ce.Error)
	UnsetPrimary(ctx context.Context, authID uint64) (a *models.Address, err *ce.Error)
}

type addressDatabase struct {
	database *db.Database
}

func NewAddressDatabase(db *db.Database) AddressDatabase {
	return &addressDatabase{database: db}
}

func (d *addressDatabase) Create(ctx context.Context, data *models.CreateAddress) (*models.Address, *ce.Error) {
	query := `
		INSERT INTO addresses (
			auth_id, label, recipient, phone, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postal_code, latitude, longitude, location
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, ST_SetSRID(ST_MakePoint($15, $14), 4326)
		)
		RETURNING
			id, label, recipient, phone, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postal_code, latitude, longitude, created_at, updated_at
	`

	var a models.Address
	err := d.database.Query(
		ctx, query,
		data.AuthID,
		data.Label,
		data.Recipient,
		data.Phone,
		data.Notes,
		data.IsPrimary,
		data.Country,
		data.Subdivision1,
		data.Subdivision2,
		data.Subdivision3,
		data.Subdivision4,
		data.Street,
		data.PostalCode,
		data.Latitude,
		data.Longitude,
	).Scan(
		&a.ID,
		&a.Label,
		&a.Recipient,
		&a.Phone,
		&a.Notes,
		&a.IsPrimary,
		&a.Country,
		&a.Subdivision1,
		&a.Subdivision2,
		&a.Subdivision3,
		&a.Subdivision4,
		&a.Street,
		&a.PostalCode,
		&a.Latitude,
		&a.Longitude,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create address: %w", err),
		)
	}

	return &a, nil
}

func (d *addressDatabase) UnsetPrimary(ctx context.Context, authID uint64) (*models.Address, *ce.Error) {
	query := `
		UPDATE
			addresses
		SET
			is_primary = FALSE,
			updated_at = NOW()
		WHERE
			auth_id = $1
			AND is_primary = TRUE
		RETURNING
			id, label, recipient, phone, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postal_code, latitude, longitude, created_at, updated_at
	`

	var a models.Address
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&a.ID,
		&a.Label,
		&a.Recipient,
		&a.Phone,
		&a.Notes,
		&a.IsPrimary,
		&a.Country,
		&a.Subdivision1,
		&a.Subdivision2,
		&a.Subdivision3,
		&a.Subdivision4,
		&a.Street,
		&a.PostalCode,
		&a.Latitude,
		&a.Longitude,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, nil
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to unset primary address: %w", err),
		)
	}

	return &a, nil
}
