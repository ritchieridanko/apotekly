package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var oAuthServiceField logger.Field = logger.NewField("service", "auth.oauth")

type OAuthClient interface {
	GetAuthURL(ctx context.Context, provider string) (url string, err *ce.Error)
	Callback(ctx context.Context, req *models.OAuthCallbackReq) (at *models.AuthToken, err *ce.Error)
}

type oAuthClient struct {
	client apis.OAuthServiceClient
}

func NewOAuthClient(c apis.OAuthServiceClient) OAuthClient {
	return &oAuthClient{client: c}
}

func (c *oAuthClient) GetAuthURL(ctx context.Context, provider string) (string, *ce.Error) {
	resp, err := c.client.GetAuthURL(
		ctx,
		&apis.GetAuthURLRequest{
			Provider: provider,
		},
	)
	if err != nil {
		return "", ce.ToError(
			err,
		).Append(
			oAuthServiceField,
		)
	}
	return resp.GetUrl(), nil
}

func (c *oAuthClient) Callback(ctx context.Context, req *models.OAuthCallbackReq) (*models.AuthToken, *ce.Error) {
	resp, err := c.client.Callback(
		ctx,
		&apis.OAuthCallbackRequest{
			Code:     req.Code,
			State:    req.State,
			Provider: req.Provider,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			oAuthServiceField,
		)
	}
	return c.toAuthToken(resp.GetAuthToken()), nil
}

func (c *oAuthClient) toAuthToken(at *apis.AuthToken) *models.AuthToken {
	if at == nil {
		return nil
	}
	return &models.AuthToken{
		AccessToken:  c.toAccessToken(at.GetAccessToken()),
		RefreshToken: c.toRefreshToken(at.GetRefreshToken()),
	}
}

func (c *oAuthClient) toAccessToken(at *apis.AccessToken) *models.AccessToken {
	if at == nil {
		return nil
	}
	return &models.AccessToken{
		Token:            at.GetToken(),
		ExpiresInSeconds: at.GetExpiresInSeconds(),
	}
}

func (c *oAuthClient) toRefreshToken(rt *apis.RefreshToken) *models.RefreshToken {
	if rt == nil {
		return nil
	}
	return &models.RefreshToken{
		Token:            rt.GetToken(),
		ExpiresInSeconds: rt.GetExpiresInSeconds(),
	}
}
