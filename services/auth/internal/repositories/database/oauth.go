package database

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type OAuthDatabase interface {
	Create(ctx context.Context, data *models.CreateOAuth) (oa *models.OAuth, err *ce.Error)
}

type oAuthDatabase struct {
	database *db.Database
}

func NewOAuthDatabase(db *db.Database) OAuthDatabase {
	return &oAuthDatabase{database: db}
}

func (d *oAuthDatabase) Create(ctx context.Context, data *models.CreateOAuth) (*models.OAuth, *ce.Error) {
	query := `
		INSERT INTO oauth (auth_id, provider, subject)
		VALUES ($1, $2, $3)
		RETURNING id, provider, subject
	`

	var oa models.OAuth
	err := d.database.Query(
		ctx, query,
		data.AuthID,
		data.Provider,
		data.Subject,
	).Scan(
		&oa.ID,
		&oa.Provider,
		&oa.Subject,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create oauth: %w", err),
		)
	}

	return &oa, nil
}
