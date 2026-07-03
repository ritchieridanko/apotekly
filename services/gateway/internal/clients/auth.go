package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var authServiceField logger.Field = logger.NewField("service", "auth")

type AuthClient interface {
	SignUp(ctx context.Context, req *models.SignUpReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	SignIn(ctx context.Context, req *models.SignInReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
}

type authClient struct {
	client apis.AuthServiceClient
}

func NewAuthClient(c apis.AuthServiceClient) AuthClient {
	return &authClient{client: c}
}

func (c *authClient) SignUp(ctx context.Context, req *models.SignUpReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	resp, err := c.client.SignUp(
		ctx,
		&apis.SignUpRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		return nil, nil, ce.ToError(
			err,
		).Append(
			authServiceField,
		)
	}
	return c.toAuth(resp.GetAuth()), c.toAuthToken(resp.GetAuthToken()), nil
}

func (c *authClient) SignIn(ctx context.Context, req *models.SignInReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	resp, err := c.client.SignIn(
		ctx,
		&apis.SignInRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		return nil, nil, ce.ToError(
			err,
		).Append(
			authServiceField,
		)
	}
	return c.toAuth(resp.GetAuth()), c.toAuthToken(resp.GetAuthToken()), nil
}

func (c *authClient) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	resp, err := c.client.IsEmailAvailable(
		ctx,
		&apis.IsEmailAvailableRequest{
			Email: email,
		},
	)
	if err != nil {
		return false, ce.ToError(
			err,
		).Append(
			authServiceField,
		)
	}
	return resp.GetIsAvailable(), nil
}

func (c *authClient) toAuth(a *apis.Auth) *models.Auth {
	if a == nil {
		return nil
	}
	return &models.Auth{
		Email:           a.GetEmail(),
		Role:            a.GetRole(),
		IsEmailVerified: a.GetIsEmailVerified(),
	}
}

func (c *authClient) toAuthToken(at *apis.AuthToken) *models.AuthToken {
	if at == nil {
		return nil
	}
	return &models.AuthToken{
		AccessToken:  c.toAccessToken(at.GetAccessToken()),
		RefreshToken: c.toRefreshToken(at.GetRefreshToken()),
	}
}

func (c *authClient) toAccessToken(at *apis.AccessToken) *models.AccessToken {
	if at == nil {
		return nil
	}
	return &models.AccessToken{
		Token:            at.GetToken(),
		ExpiresInSeconds: at.GetExpiresInSeconds(),
	}
}

func (c *authClient) toRefreshToken(rt *apis.RefreshToken) *models.RefreshToken {
	if rt == nil {
		return nil
	}
	return &models.RefreshToken{
		Token:            rt.GetToken(),
		ExpiresInSeconds: rt.GetExpiresInSeconds(),
	}
}
