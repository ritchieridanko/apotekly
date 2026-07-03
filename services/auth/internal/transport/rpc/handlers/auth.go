package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct {
	apis.UnimplementedAuthServiceServer
	au usecases.AuthUsecase
}

func NewAuthHandler(au usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{au: au}
}

func (h *AuthHandler) SignUp(ctx context.Context, req *apis.SignUpRequest) (*apis.SignUpResponse, error) {
	a, at, err := h.au.SignUp(
		ctx,
		&models.SignUpReq{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		},
	)
	if err != nil {
		return nil, err
	}
	if at == nil {
		a = nil
	}
	return &apis.SignUpResponse{
		Auth:      h.toAuth(a),
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AuthHandler) SignIn(ctx context.Context, req *apis.SignInRequest) (*apis.SignInResponse, error) {
	a, at, err := h.au.SignIn(
		ctx,
		&models.SignInReq{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.SignInResponse{
		Auth:      h.toAuth(a),
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AuthHandler) SignOut(ctx context.Context, req *apis.SignOutRequest) (*emptypb.Empty, error) {
	if err := h.au.SignOut(ctx, req.GetRefreshToken()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) IsEmailAvailable(ctx context.Context, req *apis.IsEmailAvailableRequest) (*apis.IsEmailAvailableResponse, error) {
	available, err := h.au.IsEmailAvailable(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &apis.IsEmailAvailableResponse{
		IsAvailable: available,
	}, nil
}

func (h *AuthHandler) RotateAuthToken(ctx context.Context, req *apis.RotateAuthTokenRequest) (*apis.RotateAuthTokenResponse, error) {
	at, err := h.au.RotateAuthToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	return &apis.RotateAuthTokenResponse{
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AuthHandler) toAuth(a *models.Auth) *apis.Auth {
	if a == nil {
		return nil
	}
	return &apis.Auth{
		Email:           a.Email,
		Role:            a.Role,
		IsEmailVerified: a.IsEmailVerified(),
	}
}

func (h *AuthHandler) toAuthToken(at *models.AuthToken) *apis.AuthToken {
	if at == nil {
		return nil
	}
	return &apis.AuthToken{
		AccessToken:  h.toAccessToken(at.AccessToken),
		RefreshToken: h.toRefreshToken(at.RefreshToken),
	}
}

func (h *AuthHandler) toAccessToken(at *models.AccessToken) *apis.AccessToken {
	if at == nil {
		return nil
	}
	return &apis.AccessToken{
		Token:            at.Token,
		ExpiresInSeconds: at.ExpiresInSeconds,
	}
}

func (h *AuthHandler) toRefreshToken(rt *models.RefreshToken) *apis.RefreshToken {
	if rt == nil {
		return nil
	}
	return &apis.RefreshToken{
		Token:            rt.Token,
		ExpiresInSeconds: rt.ExpiresInSeconds,
	}
}
