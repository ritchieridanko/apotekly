package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type AuthDatabase interface {
	Create(ctx context.Context, data *models.CreateAuth) (a *models.Auth, err *ce.Error)
	GetByEmail(ctx context.Context, email string) (a *models.Auth, err *ce.Error)
	GetByID(ctx context.Context, authID uint64) (a *models.Auth, err *ce.Error)
	EmailExists(ctx context.Context, email string) (exists bool, err *ce.Error)
	SetVerified(ctx context.Context, authID uint64) (a *models.Auth, err *ce.Error)
}

type authDatabase struct {
	database *db.Database
}

func NewAuthDatabase(db *db.Database) AuthDatabase {
	return &authDatabase{database: db}
}

func (d *authDatabase) Create(ctx context.Context, data *models.CreateAuth) (*models.Auth, *ce.Error) {
	query := `
		INSERT INTO auth (email, password, role, email_verified_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, role, email_verified_at
	`

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		data.Email,
		data.Password,
		data.Role,
		data.EmailVerifiedAt,
	).Scan(
		&a.ID,
		&a.Email,
		&a.Role,
		&a.EmailVerifiedAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create auth: %w", err),
		)
	}

	return &a, nil
}

func (d *authDatabase) GetByEmail(ctx context.Context, email string) (*models.Auth, *ce.Error) {
	query := `
		SELECT
			id, email, password, role, email_verified_at
		FROM
			auth
		WHERE
			email = $1
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		email,
	).Scan(
		&a.ID,
		&a.Email,
		&a.Password,
		&a.Role,
		&a.EmailVerifiedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get auth by email: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeAuthNotFound,
				ce.MsgAuthNotFound,
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

func (d *authDatabase) GetByID(ctx context.Context, authID uint64) (*models.Auth, *ce.Error) {
	query := `
		SELECT
			id, email, role, email_verified_at
		FROM
			auth
		WHERE
			id = $1
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&a.ID,
		&a.Email,
		&a.Role,
		&a.EmailVerifiedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get auth by id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeAuthNotFound,
				ce.MsgAuthNotFound,
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

func (d *authDatabase) EmailExists(ctx context.Context, email string) (bool, *ce.Error) {
	query := "SELECT 1 FROM auth WHERE email = $1 AND deleted_at IS NULL"
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var exists int
	err := d.database.Query(
		ctx, query,
		email,
	).Scan(
		&exists,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		return false, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to check if email exists: %w", err),
		)
	}

	return true, nil
}

func (d *authDatabase) SetVerified(ctx context.Context, authID uint64) (*models.Auth, *ce.Error) {
	query := `
		UPDATE
			auth
		SET
			email_verified_at = NOW(),
			updated_at = NOW()
		WHERE
			id = $1
			AND deleted_at IS NULL
		RETURNING
			id, email, role, email_verified_at
	`

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&a.ID,
		&a.Email,
		&a.Role,
		&a.EmailVerifiedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to set auth verified: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeAuthNotFound,
				ce.MsgAuthNotFound,
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
