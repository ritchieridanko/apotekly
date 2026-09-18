package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var authServiceField logger.Field = logger.NewField("service", "auth")

type AuthClient interface {
	SetRolePharmacy(ctx context.Context, authID uint64) (err *ce.Error)
}

type authClient struct {
	client apis.AuthServiceClient
}

func NewAuthClient(c apis.AuthServiceClient) AuthClient {
	return &authClient{client: c}
}

func (c *authClient) SetRolePharmacy(ctx context.Context, authID uint64) *ce.Error {
	_, err := c.client.SetRolePharmacy(
		ctx,
		&apis.SetRolePharmacyRequest{
			AuthId: authID,
		},
	)
	if err != nil {
		return ce.ToError(
			err,
		).Append(
			authServiceField,
		)
	}
	return nil
}
