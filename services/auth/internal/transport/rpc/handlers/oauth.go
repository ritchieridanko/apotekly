package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
)

type OAuthHandler struct {
	apis.UnimplementedOAuthServiceServer
	oau usecases.OAuthUsecase
}

func NewOAuthHandler(oau usecases.OAuthUsecase) *OAuthHandler {
	return &OAuthHandler{oau: oau}
}

func (h *OAuthHandler) GetAuthURL(ctx context.Context, req *apis.GetAuthURLRequest) (*apis.GetAuthURLResponse, error) {
	url, err := h.oau.GetAuthURL(ctx, req.GetProvider())
	if err != nil {
		return nil, err
	}
	return &apis.GetAuthURLResponse{Url: url}, nil
}

func (h *OAuthHandler) Callback(ctx context.Context, req *apis.OAuthCallbackRequest) (*apis.OAuthCallbackResponse, error) {
	at, err := h.oau.Callback(
		ctx,
		&models.OAuthCallbackReq{
			Code:     req.GetCode(),
			State:    req.GetState(),
			Provider: req.GetProvider(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.OAuthCallbackResponse{AuthToken: h.toAuthToken(at)}, nil
}

func (h *OAuthHandler) toAuthToken(at *models.AuthToken) *apis.AuthToken {
	if at == nil {
		return nil
	}
	return &apis.AuthToken{
		AccessToken:  h.toAccessToken(at.AccessToken),
		RefreshToken: h.toRefreshToken(at.RefreshToken),
	}
}

func (h *OAuthHandler) toAccessToken(at *models.AccessToken) *apis.AccessToken {
	if at == nil {
		return nil
	}
	return &apis.AccessToken{
		Token:            at.Token,
		ExpiresInSeconds: at.ExpiresInSeconds,
	}
}

func (h *OAuthHandler) toRefreshToken(rt *models.RefreshToken) *apis.RefreshToken {
	if rt == nil {
		return nil
	}
	return &apis.RefreshToken{
		Token:            rt.Token,
		ExpiresInSeconds: rt.ExpiresInSeconds,
	}
}
