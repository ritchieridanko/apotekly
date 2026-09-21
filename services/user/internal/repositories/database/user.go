package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
)

type UserDatabase interface {
	Create(ctx context.Context, data *models.CreateUser) (u *models.User, err *ce.Error)
	GetID(ctx context.Context, authID uint64) (userID uuid.UUID, err *ce.Error)
	GetByAuthID(ctx context.Context, authID uint64) (u *models.User, err *ce.Error)
	Update(ctx context.Context, authID uint64, data *models.UpdateUser) (u *models.User, err *ce.Error)
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

func (d *userDatabase) GetID(ctx context.Context, authID uint64) (uuid.UUID, *ce.Error) {
	query := "SELECT id FROM users WHERE auth_id = $1 AND deleted_at IS NULL"
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var userID uuid.UUID
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&userID,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get user id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return uuid.Nil, ce.NewError(
				ce.CodeUserNotFound,
				ce.MsgUserNotFound,
				wrappedErr,
			)
		}
		return uuid.Nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return userID, nil
}

func (d *userDatabase) GetByAuthID(ctx context.Context, authID uint64) (*models.User, *ce.Error) {
	query := `
		SELECT
			id, name, sex, birthdate, phone,
			profile_picture, profile_banner
		FROM
			users
		WHERE
			auth_id = $1
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var u models.User
	err := d.database.Query(
		ctx, query,
		authID,
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
		wrappedErr := fmt.Errorf("failed to get user by auth id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeUserNotFound,
				ce.MsgUserNotFound,
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

func (d *userDatabase) Update(ctx context.Context, authID uint64, data *models.UpdateUser) (*models.User, *ce.Error) {
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if data.Name != nil {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argPos))
		args = append(args, *data.Name)
		argPos++
	}
	if data.Sex != nil {
		setClauses = append(setClauses, "sex = $"+strconv.Itoa(argPos))
		args = append(args, *data.Sex)
		argPos++
	}
	if data.Birthdate != nil {
		setClauses = append(setClauses, "birthdate = $"+strconv.Itoa(argPos))
		args = append(args, *data.Birthdate)
		argPos++
	}
	if data.Phone != nil {
		setClauses = append(setClauses, "phone = $"+strconv.Itoa(argPos))
		args = append(args, *data.Phone)
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
				users
			SET
				%s
			WHERE
				auth_id = $%d
			RETURNING
				id, name, sex, birthdate, phone,
				profile_picture, profile_banner
		`,
		strings.Join(setClauses, ", "), argPos,
	)

	var u models.User
	err := d.database.Query(
		ctx, query,
		args...,
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
		wrappedErr := fmt.Errorf("failed to update user: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeUserNotFound,
				ce.MsgUserNotFound,
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
