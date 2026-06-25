package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/clients"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/dtos"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/cookie"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
)

type AuthHandler struct {
	ac        clients.AuthClient
	validator *validator.Validator
	cookie    *cookie.Cookie
}

func NewAuthHandler(ac clients.AuthClient, v *validator.Validator, c *cookie.Cookie) *AuthHandler {
	return &AuthHandler{
		ac:        ac,
		validator: v,
		cookie:    c,
	}
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	var payload dtos.SignUpRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	ip, ua := ctx.ClientIP(), ctx.Request.UserAgent()
	if ok, why := h.validator.IPAddress(ip); !ok {
		ce.NewError(ce.CodeInvalidRequestMetadata, why, nil).Bind(ctx)
		return
	}
	if ok, why := h.validator.UserAgent(ua); !ok {
		ce.NewError(ce.CodeInvalidRequestMetadata, why, nil).Bind(ctx)
		return
	}

	a, at, err := h.ac.SignUp(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyIPAddress,
			ip,
			constants.MDKeyUserAgent,
			ua,
		),
		&models.SignUpReq{
			Email:    payload.Email,
			Password: payload.Password,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}
	if at != nil && at.RefreshToken != nil {
		h.cookie.Set(
			ctx,
			constants.CookieKeyRefreshToken,
			at.RefreshToken.Token,
			"/",
			int(at.RefreshToken.ExpiresInSeconds),
		)
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusCreated,
		"Signed up successfully",
		dtos.SignUpResponse{
			Auth:        h.toAuth(a),
			AccessToken: h.toAccessToken(at),
		},
		nil,
	)
}

func (h *AuthHandler) IsEmailAvailable(ctx *gin.Context) {
	var params dtos.IsEmailAvailableRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
		return
	}

	available, err := h.ac.IsEmailAvailable(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
		),
		params.Email,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"OK",
		dtos.IsEmailAvailableResponse{
			IsAvailable: available,
		},
		nil,
	)
}

func (h *AuthHandler) toAuth(a *models.Auth) *dtos.Auth {
	if a == nil {
		return nil
	}
	return &dtos.Auth{
		Email:           a.Email,
		Role:            a.Role,
		IsEmailVerified: a.IsEmailVerified,
	}
}

func (h *AuthHandler) toAccessToken(at *models.AuthToken) *dtos.AccessToken {
	if at == nil || at.AccessToken == nil {
		return nil
	}
	return &dtos.AccessToken{
		Token:            at.AccessToken.Token,
		ExpiresInSeconds: at.AccessToken.ExpiresInSeconds,
	}
}
