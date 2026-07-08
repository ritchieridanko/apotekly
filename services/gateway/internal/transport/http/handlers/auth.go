package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var payload dtos.SignInRequest
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

	a, at, err := h.ac.SignIn(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyIPAddress,
			ip,
			constants.MDKeyUserAgent,
			ua,
		),
		&models.SignInReq{
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
		http.StatusOK,
		"Signed in successfully",
		dtos.SignInResponse{
			Auth:        h.toAuth(a),
			AccessToken: h.toAccessToken(at),
		},
		nil,
	)
}

func (h *AuthHandler) SignOut(ctx *gin.Context) {
	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	refreshToken, err := ctx.Cookie(constants.CookieKeyRefreshToken)
	if errors.Is(err, ce.ErrCookieNotFound) {
		ce.NewError(ce.CodeRefreshTokenNotFound, ce.MsgInvalidSession, err).Bind(ctx)
		return
	}
	if err != nil {
		ce.NewError(ce.CodeInternal, ce.MsgInternalServer, err).Bind(ctx)
		return
	}

	token := strings.TrimSpace(refreshToken)
	if token == "" {
		ce.NewError(
			ce.CodeRefreshTokenNotFound,
			ce.MsgInvalidSession,
			errors.New("refresh token is empty"),
		).Bind(
			ctx,
		)
		return
	}

	signOutErr := h.ac.SignOut(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
		token,
	)
	if signOutErr != nil {
		signOutErr.Bind(ctx)
		return
	}

	h.cookie.Unset(
		ctx,
		constants.CookieKeyRefreshToken,
		"/",
	)

	utils.SetHTTPResponse[any](ctx, http.StatusNoContent, "", nil, nil)
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

func (h *AuthHandler) RotateAuthToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(constants.CookieKeyRefreshToken)
	if errors.Is(err, ce.ErrCookieNotFound) {
		ce.NewError(ce.CodeRefreshTokenNotFound, ce.MsgInvalidSession, err).Bind(ctx)
		return
	}
	if err != nil {
		ce.NewError(ce.CodeInternal, ce.MsgInternalServer, err).Bind(ctx)
		return
	}

	token := strings.TrimSpace(refreshToken)
	if token == "" {
		ce.NewError(
			ce.CodeRefreshTokenNotFound,
			ce.MsgInvalidSession,
			errors.New("refresh token is empty"),
		).Bind(
			ctx,
		)
		return
	}

	at, rotateErr := h.ac.RotateAuthToken(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
		),
		token,
	)
	if rotateErr != nil {
		rotateErr.Bind(ctx)
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
		http.StatusOK,
		"Auth token rotated successfully",
		dtos.RotateAuthTokenResponse{
			AccessToken: h.toAccessToken(at),
		},
		nil,
	)
}

func (h *AuthHandler) ResendVerification(ctx *gin.Context) {
	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	email, err := h.ac.ResendVerification(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Verification resent successfully",
		dtos.ResendVerificationResponse{
			Email: email,
		},
		nil,
	)
}

func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {
	var params dtos.VerifyEmailRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
		return
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	refreshToken, err := ctx.Cookie(constants.CookieKeyRefreshToken)
	if errors.Is(err, ce.ErrCookieNotFound) {
		ce.NewError(ce.CodeRefreshTokenNotFound, ce.MsgInvalidSession, err).Bind(ctx)
		return
	}
	if err != nil {
		ce.NewError(ce.CodeInternal, ce.MsgInternalServer, err).Bind(ctx)
		return
	}

	token := strings.TrimSpace(refreshToken)
	if token == "" {
		ce.NewError(
			ce.CodeRefreshTokenNotFound,
			ce.MsgInvalidSession,
			errors.New("refresh token is empty"),
		).Bind(
			ctx,
		)
		return
	}

	a, at, verifyErr := h.ac.VerifyEmail(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
		&models.VerifyEmailReq{
			RefreshToken:      token,
			VerificationToken: params.VerificationToken,
		},
	)
	if verifyErr != nil {
		verifyErr.Bind(ctx)
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
		http.StatusOK,
		"Email verified successfully",
		dtos.VerifyEmailResponse{
			Auth:        h.toAuth(a),
			AccessToken: h.toAccessToken(at),
		},
		nil,
	)
}

func (h *AuthHandler) ChangeEmail(ctx *gin.Context) {
	var payload dtos.ChangeEmailRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	email, err := h.ac.ChangeEmail(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
		&models.ChangeEmailReq{
			Password: payload.Password,
			NewEmail: payload.NewEmail,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Email change confirmation sent successfully",
		dtos.ChangeEmailResponse{
			Email: email,
		},
		nil,
	)
}

func (h *AuthHandler) ConfirmEmailChange(ctx *gin.Context) {
	var params dtos.ConfirmEmailChangeRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
		return
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	err := h.ac.ConfirmEmailChange(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
		params.EmailChangeToken,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	h.cookie.Unset(
		ctx,
		constants.CookieKeyRefreshToken,
		"/",
	)

	utils.SetHTTPResponse[any](
		ctx,
		http.StatusOK,
		"Email changed successfully",
		nil,
		nil,
	)
}

func (h *AuthHandler) ChangePassword(ctx *gin.Context) {
	var payload dtos.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	err := h.ac.ChangePassword(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
		&models.ChangePasswordReq{
			OldPassword: payload.OldPassword,
			NewPassword: payload.NewPassword,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse[any](
		ctx,
		http.StatusOK,
		"Password changed successfully",
		nil,
		nil,
	)
}

func (h *AuthHandler) ResetPassword(ctx *gin.Context) {
	var payload dtos.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	email, err := h.ac.ResetPassword(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
		),
		payload.Email,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Password reset confirmation sent successfully",
		dtos.ResetPasswordResponse{
			Email: email,
		},
		nil,
	)
}

func (h *AuthHandler) IsPasswordResetTokenValid(ctx *gin.Context) {
	var params dtos.IsPasswordResetTokenValidRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
		return
	}

	valid, err := h.ac.IsPasswordResetTokenValid(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
		),
		params.Token,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"OK",
		dtos.IsPasswordResetTokenValidResponse{
			IsValid: valid,
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
