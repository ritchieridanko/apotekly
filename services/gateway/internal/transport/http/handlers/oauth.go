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

type OAuthHandler struct {
	clientAddr string
	oac        clients.OAuthClient
	validator  *validator.Validator
	cookie     *cookie.Cookie
}

func NewOAuthHandler(clientAddr string, oac clients.OAuthClient, v *validator.Validator, c *cookie.Cookie) *OAuthHandler {
	return &OAuthHandler{
		clientAddr: clientAddr,
		oac:        oac,
		validator:  v,
		cookie:     c,
	}
}

func (h *OAuthHandler) GetAuthURL(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if ok, why := h.validator.OAuthProvider(provider); !ok {
		ce.NewError(ce.CodeInvalidOAuthProvider, why, nil).Bind(ctx)
		return
	}

	url, err := h.oac.GetAuthURL(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
		),
		provider,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	ctx.Redirect(http.StatusFound, url)
}

func (h *OAuthHandler) Callback(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if ok, why := h.validator.OAuthProvider(provider); !ok {
		ce.NewError(ce.CodeInvalidOAuthProvider, why, nil).Bind(ctx)
		return
	}

	var params dtos.OAuthCallbackRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
		return
	}
	if params.Error != "" {
		ce.NewError(ce.CodeOAuthCancelled, params.ErrorDesc, nil).Bind(ctx)
		return
	}
	if params.Code == "" || params.State == "" {
		ce.NewError(ce.CodeInvalidParams, "Missing authorization code or state", nil).Bind(ctx)
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

	at, err := h.oac.Callback(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyIPAddress, ip,
			constants.MDKeyUserAgent, ua,
		),
		&models.OAuthCallbackReq{
			Code:     params.Code,
			State:    params.State,
			Provider: provider,
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

	ctx.Redirect(
		http.StatusTemporaryRedirect,
		h.clientAddr+"/auth/oauth-callback?token="+at.AccessToken.Token,
	)
}
