package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/models"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type PharmacyDatabase interface {
	Create(ctx context.Context, data *models.CreatePharmacy) (p *models.Pharmacy, err *ce.Error)
	GetID(ctx context.Context, authID uint64) (pharmacyID uuid.UUID, err *ce.Error)
	GetByAuthID(ctx context.Context, authID uint64) (p *models.Pharmacy, err *ce.Error)
	Update(ctx context.Context, authID uint64, data *models.UpdatePharmacy) (p *models.Pharmacy, err *ce.Error)
}

type pharmacyDatabase struct {
	database *db.Database
}

func NewPharmacyDatabase(db *db.Database) PharmacyDatabase {
	return &pharmacyDatabase{database: db}
}

func (d *pharmacyDatabase) Create(ctx context.Context, data *models.CreatePharmacy) (*models.Pharmacy, *ce.Error) {
	query := `
		INSERT INTO pharmacies (
			id, auth_id, name, legal_name, description, online_hours, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4, street,
			postal_code, latitude, longitude, location, email, phone, website, whatsapp
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			ST_SetSRID(ST_MakePoint($15, $14), 4326), $16, $17, $18, $19
		)
		RETURNING
			id, name, legal_name, description, status, online_hours, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4, street,
			postal_code, latitude, longitude, email, phone, website, whatsapp,
			profile_picture, profile_banner, verified_at, created_at, updated_at
	`

	var p models.Pharmacy
	err := d.database.Query(
		ctx, query,
		data.ID,
		data.AuthID,
		data.Name,
		data.LegalName,
		data.Description,
		data.OnlineHours,
		data.Country,
		data.Subdivision1,
		data.Subdivision2,
		data.Subdivision3,
		data.Subdivision4,
		data.Street,
		data.PostalCode,
		data.Latitude,
		data.Longitude,
		data.Email,
		data.Phone,
		data.Website,
		data.Whatsapp,
	).Scan(
		&p.ID,
		&p.Name,
		&p.LegalName,
		&p.Description,
		&p.Status,
		&p.OnlineHours,
		&p.Country,
		&p.Subdivision1,
		&p.Subdivision2,
		&p.Subdivision3,
		&p.Subdivision4,
		&p.Street,
		&p.PostalCode,
		&p.Latitude,
		&p.Longitude,
		&p.Email,
		&p.Phone,
		&p.Website,
		&p.Whatsapp,
		&p.ProfilePicture,
		&p.ProfileBanner,
		&p.VerifiedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create pharmacy: %w", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Auth ID already exists
			return nil, ce.NewError(
				ce.CodePharmacyAlreadyExists,
				ce.MsgPharmacyAlreadyExists,
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

func (d *pharmacyDatabase) GetID(ctx context.Context, authID uint64) (uuid.UUID, *ce.Error) {
	query := "SELECT id FROM pharmacies WHERE auth_id = $1 AND deleted_at IS NULL"
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var pharmacyID uuid.UUID
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&pharmacyID,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get pharmacy id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return uuid.Nil, ce.NewError(
				ce.CodePharmacyNotFound,
				ce.MsgPharmacyNotFound,
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

func (d *pharmacyDatabase) GetByAuthID(ctx context.Context, authID uint64) (*models.Pharmacy, *ce.Error) {
	query := `
		SELECT
			id, name, legal_name, description, status, online_hours, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4, street,
			postal_code, latitude, longitude, email, phone, website, whatsapp,
			profile_picture, profile_banner, verified_at, created_at, updated_at
		FROM
			pharmacies
		WHERE
			auth_id = $1
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var p models.Pharmacy
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&p.ID,
		&p.Name,
		&p.LegalName,
		&p.Description,
		&p.Status,
		&p.OnlineHours,
		&p.Country,
		&p.Subdivision1,
		&p.Subdivision2,
		&p.Subdivision3,
		&p.Subdivision4,
		&p.Street,
		&p.PostalCode,
		&p.Latitude,
		&p.Longitude,
		&p.Email,
		&p.Phone,
		&p.Website,
		&p.Whatsapp,
		&p.ProfilePicture,
		&p.ProfileBanner,
		&p.VerifiedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get pharmacy by auth id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodePharmacyNotFound,
				ce.MsgPharmacyNotFound,
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

func (d *pharmacyDatabase) Update(ctx context.Context, authID uint64, data *models.UpdatePharmacy) (*models.Pharmacy, *ce.Error) {
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if data.Name != nil {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argPos))
		args = append(args, *data.Name)
		argPos++
	}
	if data.LegalName != nil {
		setClauses = append(setClauses, "legal_name = $"+strconv.Itoa(argPos))
		args = append(args, *data.LegalName)
		argPos++
	}
	if data.Description != nil {
		setClauses = append(setClauses, "description = $"+strconv.Itoa(argPos))
		args = append(args, *data.Description)
		argPos++
	}
	if data.OnlineHours != nil {
		setClauses = append(setClauses, "online_hours = $"+strconv.Itoa(argPos))
		args = append(args, *data.OnlineHours)
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
	if data.Email != nil {
		setClauses = append(setClauses, "email = $"+strconv.Itoa(argPos))
		args = append(args, *data.Email)
		argPos++
	}
	if data.Phone != nil {
		setClauses = append(setClauses, "phone = $"+strconv.Itoa(argPos))
		args = append(args, *data.Phone)
		argPos++
	}
	if data.Website != nil {
		setClauses = append(setClauses, "website = $"+strconv.Itoa(argPos))
		args = append(args, *data.Website)
		argPos++
	}
	if data.Whatsapp != nil {
		setClauses = append(setClauses, "whatsapp = $"+strconv.Itoa(argPos))
		args = append(args, *data.Whatsapp)
		argPos++
	}
	if len(setClauses) == 0 {
		return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, nil)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, authID)
	query := fmt.Sprintf(
		`
			UPDATE
				pharmacies
			SET
				%s
			WHERE
				auth_id = $%d
			RETURNING
				id, name, legal_name, description, status, online_hours, country,
				subdivision_1, subdivision_2, subdivision_3, subdivision_4, street,
				postal_code, latitude, longitude, email, phone, website, whatsapp,
				profile_picture, profile_banner, verified_at, created_at, updated_at
		`,
		strings.Join(setClauses, ", "), argPos,
	)

	var p models.Pharmacy
	err := d.database.Query(
		ctx, query,
		args...,
	).Scan(
		&p.ID,
		&p.Name,
		&p.LegalName,
		&p.Description,
		&p.Status,
		&p.OnlineHours,
		&p.Country,
		&p.Subdivision1,
		&p.Subdivision2,
		&p.Subdivision3,
		&p.Subdivision4,
		&p.Street,
		&p.PostalCode,
		&p.Latitude,
		&p.Longitude,
		&p.Email,
		&p.Phone,
		&p.Website,
		&p.Whatsapp,
		&p.ProfilePicture,
		&p.ProfileBanner,
		&p.VerifiedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update pharmacy: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodePharmacyNotFound,
				ce.MsgPharmacyNotFound,
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
