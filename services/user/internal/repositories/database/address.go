package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
)

type AddressDatabase interface {
	Create(ctx context.Context, data *models.CreateAddress) (a *models.Address, err *ce.Error)
	GetAll(ctx context.Context, params *models.GetAllAddresses) (as []models.Address, total int64, err *ce.Error)
	Update(ctx context.Context, params *models.UpdateAddressP, data *models.UpdateAddressD) (a *models.Address, err *ce.Error)
	Delete(ctx context.Context, params *models.DeleteAddress) (err *ce.Error)
	SetPrimary(ctx context.Context, params *models.SetPrimaryAddress) (a *models.Address, err *ce.Error)
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

func (d *addressDatabase) GetAll(ctx context.Context, params *models.GetAllAddresses) ([]models.Address, int64, *ce.Error) {
	query := `
		SELECT
			id, label, recipient, phone, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postal_code, latitude, longitude, created_at, updated_at,

			COUNT(*) OVER() AS total
		FROM
			addresses
		WHERE
			auth_id = $1
		ORDER BY
			is_primary DESC,
			updated_at DESC,
			id DESC
		LIMIT
			$2
		OFFSET
			$3
	`

	rows, err := d.database.QueryAll(
		ctx, query,
		params.AuthID,
		params.PageSize,
		params.Offset(),
	)
	if err != nil {
		return nil, 0, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all addresses: %w", err),
		)
	}
	defer rows.Close()

	var total int64
	as := make([]models.Address, 0, params.PageSize)

	for rows.Next() {
		var a models.Address
		err := rows.Scan(
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
			&total,
		)
		if err != nil {
			return nil, 0, ce.NewError(
				ce.CodeDBQueryExec,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get all addresses: %w", err),
			)
		}

		as = append(as, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all addresses: %w", err),
		)
	}
	if len(as) == 0 {
		err := d.database.Query(
			ctx,
			"SELECT COUNT(*) FROM addresses WHERE auth_id = $1",
			params.AuthID,
		).Scan(
			&total,
		)
		if err != nil {
			return nil, 0, ce.NewError(
				ce.CodeDBQueryExec,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get all addresses: %w", err),
			)
		}
	}

	return as, total, nil
}

func (d *addressDatabase) Update(ctx context.Context, params *models.UpdateAddressP, data *models.UpdateAddressD) (*models.Address, *ce.Error) {
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if data.Label != nil {
		setClauses = append(setClauses, "label = $"+strconv.Itoa(argPos))
		args = append(args, *data.Label)
		argPos++
	}
	if data.Recipient != nil {
		setClauses = append(setClauses, "recipient = $"+strconv.Itoa(argPos))
		args = append(args, *data.Recipient)
		argPos++
	}
	if data.Phone != nil {
		setClauses = append(setClauses, "phone = $"+strconv.Itoa(argPos))
		args = append(args, *data.Phone)
		argPos++
	}
	if data.Notes != nil {
		setClauses = append(setClauses, "notes = $"+strconv.Itoa(argPos))
		args = append(args, *data.Notes)
		argPos++
	}
	if data.Country != nil {
		setClauses = append(setClauses, "country = $"+strconv.Itoa(argPos))
		args = append(args, *data.Country)
		argPos++
	}
	if data.Subdivision1 != nil {
		setClauses = append(setClauses, "subdivision_1 = $"+strconv.Itoa(argPos))
		args = append(args, *data.Subdivision1)
		argPos++
	}
	if data.Subdivision2 != nil {
		setClauses = append(setClauses, "subdivision_2 = $"+strconv.Itoa(argPos))
		args = append(args, *data.Subdivision2)
		argPos++
	}
	if data.Subdivision3 != nil {
		setClauses = append(setClauses, "subdivision_3 = $"+strconv.Itoa(argPos))
		args = append(args, *data.Subdivision3)
		argPos++
	}
	if data.Subdivision4 != nil {
		setClauses = append(setClauses, "subdivision_4 = $"+strconv.Itoa(argPos))
		args = append(args, *data.Subdivision4)
		argPos++
	}
	if data.Street != nil {
		setClauses = append(setClauses, "street = $"+strconv.Itoa(argPos))
		args = append(args, *data.Street)
		argPos++
	}
	if data.PostalCode != nil {
		setClauses = append(setClauses, "postal_code = $"+strconv.Itoa(argPos))
		args = append(args, *data.PostalCode)
		argPos++
	}
	if data.Latitude != nil && data.Longitude != nil {
		setClauses = append(
			setClauses,
			"latitude = $"+strconv.Itoa(argPos),
			"longitude = $"+strconv.Itoa(argPos+1),
			"location = ST_SetSRID(ST_MakePoint($"+strconv.Itoa(argPos+1)+", $"+strconv.Itoa(argPos)+"), 4326)",
		)
		args = append(args, *data.Latitude, *data.Longitude)
		argPos += 2
	}
	if len(setClauses) == 0 {
		return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, nil)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, params.AddressID, params.AuthID)
	query := fmt.Sprintf(
		`
			UPDATE
				addresses
			SET
				%s
			WHERE
				id = $%d
				AND auth_id = $%d
			RETURNING
				id, label, recipient, phone, notes, is_primary, country,
				subdivision_1, subdivision_2, subdivision_3, subdivision_4,
				street, postal_code, latitude, longitude, created_at, updated_at
		`,
		strings.Join(setClauses, ", "), argPos, argPos+1,
	)

	var a models.Address
	err := d.database.Query(
		ctx, query,
		args...,
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
		wrappedErr := fmt.Errorf("failed to update address: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeAddressNotFound,
				ce.MsgAddressNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &a, nil
}

func (d *addressDatabase) Delete(ctx context.Context, params *models.DeleteAddress) *ce.Error {
	query := "DELETE FROM addresses WHERE id = $1 AND auth_id = $2"

	err := d.database.Execute(
		ctx, query,
		params.AddressID,
		params.AuthID,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete address: %w", err)
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(
				ce.CodeAddressNotFound,
				ce.MsgAddressNotFound,
				wrappedErr,
			)
		}
		return ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return nil
}

func (d *addressDatabase) SetPrimary(ctx context.Context, params *models.SetPrimaryAddress) (*models.Address, *ce.Error) {
	query := `
		UPDATE
			addresses
		SET
			is_primary = TRUE,
			updated_at = NOW()
		WHERE
			id = $1
			AND auth_id = $2
		RETURNING
			id, label, recipient, phone, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postal_code, latitude, longitude, created_at, updated_at
	`

	var a models.Address
	err := d.database.Query(
		ctx, query,
		params.AddressID,
		params.AuthID,
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
		wrappedErr := fmt.Errorf("failed to set primary address: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeAddressNotFound,
				ce.MsgAddressNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
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
