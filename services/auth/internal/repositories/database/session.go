package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type SessionDatabase interface {
	Create(ctx context.Context, data *models.CreateSession) (err *ce.Error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (s *models.Session, err *ce.Error)
	Revoke(ctx context.Context, params *models.RevokeSession) (s *models.Session, err *ce.Error)
	RevokeActive(ctx context.Context, params *models.RevokeActiveSession) (sessionID uint64, err *ce.Error)
}

type sessionDatabase struct {
	database *db.Database
}

func NewSessionDatabase(db *db.Database) SessionDatabase {
	return &sessionDatabase{database: db}
}

func (d *sessionDatabase) Create(ctx context.Context, data *models.CreateSession) *ce.Error {
	query := `
		INSERT INTO sessions (
			parent_id, auth_id, refresh_token,
			ip_address, user_agent, expires_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	err := d.database.Execute(
		ctx, query,
		data.ParentID,
		data.AuthID,
		data.RefreshToken,
		data.IPAddress,
		data.UserAgent,
		data.ExpiresAt,
	)
	if err != nil {
		return ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create session: %w", err),
		)
	}

	return nil
}

func (d *sessionDatabase) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, *ce.Error) {
	query := `
		SELECT
			id, auth_id, refresh_token,
			ip_address, user_agent, expires_at
		FROM
			sessions
		WHERE
			refresh_token = $1
			AND revoked_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var s models.Session
	err := d.database.Query(
		ctx, query,
		refreshToken,
	).Scan(
		&s.ID,
		&s.AuthID,
		&s.RefreshToken,
		&s.IPAddress,
		&s.UserAgent,
		&s.ExpiresAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get session by refresh token: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeSessionNotFound,
				ce.MsgSessionNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &s, nil
}

func (d *sessionDatabase) Revoke(ctx context.Context, params *models.RevokeSession) (*models.Session, *ce.Error) {
	query := `
		UPDATE
			sessions
		SET
			revoked_at = NOW()
		WHERE
			refresh_token = $1
			AND revoked_at IS NULL
			AND expires_at >= $2
		RETURNING
			id, auth_id, refresh_token,
			ip_address, user_agent
	`

	var s models.Session
	err := d.database.Query(
		ctx, query,
		params.RefreshToken,
		params.ExpiresAt,
	).Scan(
		&s.ID,
		&s.AuthID,
		&s.RefreshToken,
		&s.IPAddress,
		&s.UserAgent,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to revoke session: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeSessionNotFound,
				ce.MsgSessionNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &s, nil
}

func (d *sessionDatabase) RevokeActive(ctx context.Context, params *models.RevokeActiveSession) (uint64, *ce.Error) {
	query := `
		UPDATE
			sessions
		SET
			revoked_at = NOW()
		WHERE
			auth_id = $1
			AND revoked_at IS NULL
			AND ip_address = $2
			AND user_agent = $3
			AND expires_at >= $4
		RETURNING
			id
	`

	var sessionID uint64
	err := d.database.Query(
		ctx, query,
		params.AuthID,
		params.IPAddress,
		params.UserAgent,
		params.ExpiresAt,
	).Scan(
		&sessionID,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, nil
		}
		return 0, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to revoke active session: %w", err),
		)
	}

	return sessionID, nil
}
