package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var userServiceField logger.Field = logger.NewField("service", "user")

type UserClient interface {
	CreateUser(ctx context.Context, req *models.CreateUserReq) (u *models.User, err *ce.Error)
}

type userClient struct {
	client apis.UserServiceClient
}

func NewUserClient(c apis.UserServiceClient) UserClient {
	return &userClient{client: c}
}

func (c *userClient) CreateUser(ctx context.Context, req *models.CreateUserReq) (*models.User, *ce.Error) {
	resp, err := c.client.CreateUser(
		ctx,
		&apis.CreateUserRequest{
			Name:      req.Name,
			Sex:       req.Sex,
			Birthdate: utils.ToTimestamp(req.Birthdate),
			Phone:     req.Phone,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			userServiceField,
		)
	}
	return c.toUser(resp.GetUser()), nil
}

func (c *userClient) toUser(u *apis.User) *models.User {
	if u == nil {
		return nil
	}
	return &models.User{
		ID:             utils.ToUUID(u.GetId()),
		Name:           u.GetName(),
		Sex:            u.Sex,
		Birthdate:      utils.ToTime(u.GetBirthdate()),
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		ProfileBanner:  u.ProfileBanner,
	}
}
