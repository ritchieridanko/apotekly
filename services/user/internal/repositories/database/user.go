package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
)

type UserDatabase interface {
	Create(ctx context.Context, data *models.CreateUser) (u *models.User, err *ce.Error)
}

type userDatabase struct {
	database *db.Database
}

func NewUserDatabase(db *db.Database) UserDatabase {
	return &userDatabase{database: db}
}

func (d *userDatabase) Create(ctx context.Context, data *models.CreateUser) (*models.User, *ce.Error) {
	query := `
		INSERT INTO users (
			id, auth_id, name, sex, birthdate, phone
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING
			id, name, sex, birthdate, phone,
			profile_picture, profile_banner
	`

	var u models.User
	err := d.database.Query(
		ctx, query,
		data.ID,
		data.AuthID,
		data.Name,
		data.Sex,
		data.Birthdate,
		data.Phone,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Sex,
		&u.Birthdate,
		&u.Phone,
		&u.ProfilePicture,
		&u.ProfileBanner,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create user: %w", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Auth ID already exists
			return nil, ce.NewError(
				ce.CodeUserAlreadyExists,
				ce.MsgUserAlreadyExists,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &u, nil
}
