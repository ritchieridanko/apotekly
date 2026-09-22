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
)

type PharmacyHandler struct {
	phc    clients.PharmacyClient
	ac     clients.AuthClient
	cookie *cookie.Cookie
}

func NewPharmacyHandler(phc clients.PharmacyClient, ac clients.AuthClient, c *cookie.Cookie) *PharmacyHandler {
	return &PharmacyHandler{phc: phc, ac: ac, cookie: c}
}

func (h *PharmacyHandler) CreatePharmacy(ctx *gin.Context) {
	var payload dtos.CreatePharmacyRequest
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

	// Pharmacy Creation
	p, err := h.phc.CreatePharmacy(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole,
			authCtx.Role,
			constants.MDKeyIsVerified,
			strconv.FormatBool(authCtx.IsEmailVerified),
		),
		&models.CreatePharmacyReq{
			Name:         payload.Name,
			LegalName:    payload.LegalName,
			Description:  payload.Description,
			OnlineHours:  payload.OnlineHours,
			Country:      payload.Country,
			Subdivision1: payload.Subdivision1,
			Subdivision2: payload.Subdivision2,
			Subdivision3: payload.Subdivision3,
			Subdivision4: payload.Subdivision4,
			Street:       payload.Street,
			PostalCode:   payload.PostalCode,
			Latitude:     payload.Latitude,
			Longitude:    payload.Longitude,
			Email:        payload.Email,
			Phone:        payload.Phone,
			Website:      payload.Website,
			Whatsapp:     payload.Whatsapp,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	// Auth Token Rotation
	refreshToken, getErr := ctx.Cookie(constants.CookieKeyRefreshToken)
	if getErr == nil {
		token := strings.TrimSpace(refreshToken)
		if token != "" {
			at, err := h.ac.RotateAuthToken(
				utils.CtxWithMetadata(
					ctx.Request.Context(),
				),
				token,
			)
			if err == nil {
				if at != nil && at.RefreshToken != nil {
					var duration int
					if payload.RememberMe {
						duration = int(at.RefreshToken.ExpiresInSeconds)
					}

					h.cookie.Set(
						ctx,
						constants.CookieKeyRefreshToken,
						at.RefreshToken.Token,
						"/",
						duration,
					)
				}

				utils.SetHTTPResponse(
					ctx,
					http.StatusCreated,
					"Pharmacy created successfully",
					dtos.CreatePharmacyResponse{
						Pharmacy:    h.toPharmacy(p),
						AccessToken: h.toAccessToken(at),
					},
					nil,
				)
				return
			}
		}
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusCreated,
		"Pharmacy created successfully",
		dtos.CreatePharmacyResponse{Pharmacy: h.toPharmacy(p)},
		nil,
	)
}

func (h *PharmacyHandler) GetMe(ctx *gin.Context) {
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

	p, err := h.phc.GetMe(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole,
			authCtx.Role,
		),
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Pharmacy retrieved successfully",
		dtos.PharmacyGetMeResponse{Pharmacy: h.toPharmacy(p)},
		nil,
	)
}

func (h *PharmacyHandler) UpdatePharmacy(ctx *gin.Context) {
	var payload dtos.UpdatePharmacyRequest
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

	p, err := h.phc.UpdatePharmacy(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole,
			authCtx.Role,
		),
		&models.UpdatePharmacyReq{
			Name:         payload.Name,
			LegalName:    payload.LegalName,
			Description:  payload.Description,
			OnlineHours:  payload.OnlineHours,
			Country:      payload.Country,
			Subdivision1: payload.Subdivision1,
			Subdivision2: payload.Subdivision2,
			Subdivision3: payload.Subdivision3,
			Subdivision4: payload.Subdivision4,
			Street:       payload.Street,
			PostalCode:   payload.PostalCode,
			Latitude:     payload.Latitude,
			Longitude:    payload.Longitude,
			Email:        payload.Email,
			Phone:        payload.Phone,
			Website:      payload.Website,
			Whatsapp:     payload.Whatsapp,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Pharmacy updated successfully",
		dtos.UpdatePharmacyResponse{Pharmacy: h.toPharmacy(p)},
		nil,
	)
}

func (h *PharmacyHandler) UpdateProfilePicture(ctx *gin.Context) {
	var payload dtos.PharmacyUpdateProfilePictureRequest
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
	if authCtx.PharmacyID == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("pharmacy_id missing from auth context"),
		).Bind(
			ctx,
		)
		return
	}

	p, err := h.phc.UpdateProfilePicture(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole,
			authCtx.Role,
			constants.MDKeyPharmacyID,
			authCtx.PharmacyID.String(),
		),
		payload.ProfilePictureURL,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Profile picture updated successfully",
		dtos.PharmacyUpdateProfilePictureResponse{Pharmacy: h.toPharmacy(p)},
		nil,
	)
}

func (h *PharmacyHandler) UpdateProfileBanner(ctx *gin.Context) {
	var payload dtos.PharmacyUpdateProfileBannerRequest
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
	if authCtx.PharmacyID == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("pharmacy_id missing from auth context"),
		).Bind(
			ctx,
		)
		return
	}

	p, err := h.phc.UpdateProfileBanner(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole,
			authCtx.Role,
			constants.MDKeyPharmacyID,
			authCtx.PharmacyID.String(),
		),
		payload.ProfileBannerURL,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Profile banner updated successfully",
		dtos.PharmacyUpdateProfileBannerResponse{Pharmacy: h.toPharmacy(p)},
		nil,
	)
}

func (h *PharmacyHandler) toPharmacy(p *models.Pharmacy) *dtos.Pharmacy {
	if p == nil {
		return nil
	}
	return &dtos.Pharmacy{
		ID:             p.ID.String(),
		Name:           p.Name,
		LegalName:      p.LegalName,
		Description:    p.Description,
		Status:         p.Status,
		OnlineHours:    p.OnlineHours,
		Country:        p.Country,
		Subdivision1:   p.Subdivision1,
		Subdivision2:   p.Subdivision2,
		Subdivision3:   p.Subdivision3,
		Subdivision4:   p.Subdivision4,
		Street:         p.Street,
		PostalCode:     p.PostalCode,
		Latitude:       p.Latitude,
		Longitude:      p.Longitude,
		Email:          p.Email,
		Phone:          p.Phone,
		Website:        p.Website,
		Whatsapp:       p.Whatsapp,
		ProfilePicture: p.ProfilePicture,
		ProfileBanner:  p.ProfileBanner,
		VerifiedAt:     p.VerifiedAt,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func (h *PharmacyHandler) toAccessToken(at *models.AuthToken) *dtos.AccessToken {
	if at == nil || at.AccessToken == nil {
		return nil
	}
	return &dtos.AccessToken{
		Token:            at.AccessToken.Token,
		ExpiresInSeconds: at.AccessToken.ExpiresInSeconds,
	}
}
