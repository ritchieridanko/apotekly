package clients

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var userServiceField logger.Field = logger.NewField("service", "user")

type UserClient interface {
	GetID(ctx context.Context, authID uint64) (userID uuid.UUID, err *ce.Error)
}

type userClient struct {
	client apis.UserServiceClient
}

func NewUserClient(c apis.UserServiceClient) UserClient {
	return &userClient{client: c}
}

func (c *userClient) GetID(ctx context.Context, authID uint64) (uuid.UUID, *ce.Error) {
	resp, err := c.client.GetID(
		ctx,
		&apis.UserGetIDRequest{
			AuthId: authID,
		},
	)
	if err != nil {
		return uuid.Nil, ce.ToError(
			err,
		).Append(
			userServiceField,
		)
	}
	return utils.ToUUID(resp.GetUserId()), nil
}
