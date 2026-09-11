package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/models"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type PharmacyDatabase interface {
	Create(ctx context.Context, data *models.CreatePharmacy) (p *models.Pharmacy, err *ce.Error)
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
