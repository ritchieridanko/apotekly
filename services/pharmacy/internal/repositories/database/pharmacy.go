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
	GetByID(ctx context.Context, pharmacyID uuid.UUID) (p *models.Pharmacy, err *ce.Error)
	GetByAuthID(ctx context.Context, authID uint64) (p *models.Pharmacy, err *ce.Error)
	GetAll(ctx context.Context, params *models.GetAllPharmacies) (pss []models.PharmacySummary, total int64, err *ce.Error)
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

func (d *pharmacyDatabase) GetByID(ctx context.Context, pharmacyID uuid.UUID) (*models.Pharmacy, *ce.Error) {
	query := `
		SELECT
			id, name, legal_name, description, status, online_hours, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4, street,
			postal_code, latitude, longitude, email, phone, website, whatsapp,
			profile_picture, profile_banner, verified_at, created_at, updated_at
		FROM
			pharmacies
		WHERE
			id = $1
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var p models.Pharmacy
	err := d.database.Query(
		ctx, query,
		pharmacyID,
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
		wrappedErr := fmt.Errorf("failed to get pharmacy by id: %w", err)
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

func (d *pharmacyDatabase) GetAll(ctx context.Context, params *models.GetAllPharmacies) ([]models.PharmacySummary, int64, *ce.Error) {
	query := ""
	args := []any{}
	argPos := 1
	hasLocation := params.RequireLocation()
	if hasLocation {
		query += "WITH user AS (SELECT ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography AS location) "
		args = append(args, *params.Latitude, *params.Longitude)
		argPos += 2
	}

	query += `
		SELECT
			p.id, p.name, p.legal_name, p.online_hours, p.profile_picture,
			p.created_at, p.updated_at,
	`
	if hasLocation {
		query += "ST_Distance(p.location, u.location) AS distance_m,"
	} else {
		query += "NULL::DOUBLE PRECISION AS distance_m,"
	}

	query += "COUNT(*) OVER() AS total FROM pharmacies p"
	if hasLocation {
		query += " JOIN user u ON TRUE"
	}

	query += " WHERE p.deleted_at IS NULL AND p.is_active = TRUE"
	searchArgPos := -1
	if params.Search != nil {
		query += " AND (p.name <% $" + strconv.Itoa(argPos) + " OR p.legal_name <% $" + strconv.Itoa(argPos) + ")"
		args = append(args, *params.Search)
		searchArgPos = argPos
		argPos++
	}
	if params.RadiusM != nil {
		query += " AND ST_DWithin(p.location, u.location, $" + strconv.Itoa(argPos) + ")"
		args = append(args, *params.RadiusM)
		argPos++
	}

	sortClauses := []string{}
	if params.Search != nil {
		sortClauses = append(
			sortClauses,
			"GREATEST(word_similarity(p.name, $"+strconv.Itoa(searchArgPos)+"), word_similarity(p.legal_name, $"+strconv.Itoa(searchArgPos)+")) DESC",
		)
	}
	if params.ByLocation != nil {
		if params.ByLocation.IsAsc {
			sortClauses = append(sortClauses, "p.location <-> u.location ASC")
		} else {
			sortClauses = append(sortClauses, "p.location <-> u.location DESC")
		}
	}
	if params.ByCreatedAt != nil {
		if params.ByCreatedAt.IsAsc {
			sortClauses = append(sortClauses, "p.created_at ASC")
		} else {
			sortClauses = append(sortClauses, "p.created_at DESC")
		}
	}
	if params.ByUpdatedAt != nil {
		if params.ByUpdatedAt.IsAsc {
			sortClauses = append(sortClauses, "p.updated_at ASC")
		} else {
			sortClauses = append(sortClauses, "p.updated_at DESC")
		}
	}
	if len(sortClauses) > 0 {
		query += " ORDER BY " + strings.Join(sortClauses, ", ")
	}

	query += " LIMIT $" + strconv.Itoa(argPos) + " OFFSET $" + strconv.Itoa(argPos+1)
	args = append(args, params.PageSize, params.Offset())

	rows, err := d.database.QueryAll(
		ctx, query,
		args...,
	)
	if err != nil {
		return nil, 0, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all pharmacies: %w", err),
		)
	}
	defer rows.Close()

	var total int64
	pss := make([]models.PharmacySummary, 0, params.PageSize)

	for rows.Next() {
		var ps models.PharmacySummary
		err := rows.Scan(
			&ps.ID,
			&ps.Name,
			&ps.LegalName,
			&ps.OnlineHours,
			&ps.ProfilePicture,
			&ps.CreatedAt,
			&ps.UpdatedAt,
			&ps.DistanceM,
			&total,
		)
		if err != nil {
			return nil, 0, ce.NewError(
				ce.CodeDBQueryExec,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get all pharmacies: %w", err),
			)
		}

		pss = append(pss, ps)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all pharmacies: %w", err),
		)
	}
	if len(pss) == 0 {
		countQuery := ""
		countArgs := []any{}
		countArgPos := 1
		if params.RadiusM != nil {
			countQuery += "WITH user AS (SELECT ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography AS location) "
			countArgs = append(countArgs, *params.Latitude, *params.Longitude)
			countArgPos += 2
		}

		countQuery += "SELECT COUNT(*) FROM pharmacies p"
		if params.RadiusM != nil {
			countQuery += " JOIN user u ON TRUE"
		}

		countQuery += " WHERE p.deleted_at IS NULL AND p.is_active = TRUE"
		if params.Search != nil {
			countQuery += " AND (p.name <% $" + strconv.Itoa(countArgPos) + " OR p.legal_name <% $" + strconv.Itoa(countArgPos) + ")"
			countArgs = append(countArgs, *params.Search)
			countArgPos++
		}
		if params.RadiusM != nil {
			countQuery += " AND ST_DWithin(p.location, u.location, $" + strconv.Itoa(countArgPos) + ")"
			countArgs = append(countArgs, *params.RadiusM)
			countArgPos++
		}

		err := d.database.Query(
			ctx, countQuery,
			countArgs...,
		).Scan(
			&total,
		)
		if err != nil {
			return nil, 0, ce.NewError(
				ce.CodeDBQueryExec,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get total pharmacies count: %w", err),
			)
		}
	}

	return pss, total, nil
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
	if data.ProfilePicture != nil {
		setClauses = append(setClauses, "profile_picture = $"+strconv.Itoa(argPos))
		args = append(args, *data.ProfilePicture)
		argPos++
	}
	if data.ProfileBanner != nil {
		setClauses = append(setClauses, "profile_banner = $"+strconv.Itoa(argPos))
		args = append(args, *data.ProfileBanner)
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
