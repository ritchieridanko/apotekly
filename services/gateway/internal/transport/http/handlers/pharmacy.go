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
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/cookie"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
)

type PharmacyHandler struct {
	phc       clients.PharmacyClient
	ac        clients.AuthClient
	validator *validator.Validator
	cookie    *cookie.Cookie
	logger    *logger.Logger
}

func NewPharmacyHandler(
	phc clients.PharmacyClient,
	ac clients.AuthClient,
	v *validator.Validator,
	c *cookie.Cookie,
	l *logger.Logger,
) *PharmacyHandler {
	return &PharmacyHandler{
		phc:       phc,
		ac:        ac,
		validator: v,
		cookie:    c,
		logger:    l,
	}
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

			h.logger.Warn(
				ctx,
				"created pharmacy. failed to rotate auth token",
				err.Append(
					logger.NewField("auth_id", authCtx.AuthID),
					logger.NewField("pharmacy_id", p.ID.String()),
					logger.NewField("error_code", err.Code()),
					logger.NewField("error", err.Unwrap()),
				).Fields()...,
			)
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

func (h *PharmacyHandler) GetPharmacyByID(ctx *gin.Context) {
	p, err := h.phc.GetPharmacyByID(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
		),
		ctx.Param("pharmacy_id"),
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Pharmacy retrieved successfully",
		dtos.GetPharmacyByIDResponse{Pharmacy: h.toPharmacy(p)},
		nil,
	)
}

func (h *PharmacyHandler) GetAllPharmacies(ctx *gin.Context) {
	var params dtos.GetAllPharmaciesRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
		return
	}

	params.SortLocation = utils.ToLowerPtr(utils.TrimSpacePtr(params.SortLocation))
	params.SortCreatedAt = utils.ToLowerPtr(utils.TrimSpacePtr(params.SortCreatedAt))
	params.SortUpdatedAt = utils.ToLowerPtr(utils.TrimSpacePtr(params.SortUpdatedAt))

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = constants.PageDefaultSizePharmacy
	}
	if params.PageSize > constants.PageMaxSizePharmacy {
		params.PageSize = constants.PageMaxSizePharmacy
	}
	if params.SortLocation != nil {
		if ok, why := h.validator.Sorter(*params.SortLocation, "Location sorter"); !ok {
			ce.NewError(ce.CodeInvalidParams, why, nil).Bind(ctx)
			return
		}
	}
	if params.SortCreatedAt != nil {
		if ok, why := h.validator.Sorter(*params.SortCreatedAt, "Creation time sorter"); !ok {
			ce.NewError(ce.CodeInvalidParams, why, nil).Bind(ctx)
			return
		}
	}
	if params.SortUpdatedAt != nil {
		if ok, why := h.validator.Sorter(*params.SortUpdatedAt, "Update time sorter"); !ok {
			ce.NewError(ce.CodeInvalidParams, why, nil).Bind(ctx)
			return
		}
	}

	requestCtx := ctx.Request.Context()
	if authCtx := utils.CtxAuth(requestCtx); authCtx != nil {
		requestCtx = utils.CtxWithMetadata(
			requestCtx,
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		)
	}

	pss, total, err := h.phc.GetAllPharmacies(
		requestCtx,
		&models.GetAllPharmaciesReq{
			Search:    params.Search,
			RadiusM:   params.RadiusM,
			Latitude:  params.Latitude,
			Longitude: params.Longitude,

			ByLocation:  params.SortLocation,
			ByCreatedAt: params.SortCreatedAt,
			ByUpdatedAt: params.SortUpdatedAt,

			Page:     int32(params.Page),
			PageSize: int32(params.PageSize),
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	pharmacies := make([]dtos.PharmacySummary, 0, len(pss))
	for _, ps := range pss {
		pharmacies = append(pharmacies, *h.toPharmacySummary(&ps))
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Pharmacies retrieved successfully",
		dtos.GetAllPharmaciesResponse{
			Pharmacies: pharmacies,
		},
		&utils.ResponseMetadata{
			Page:     params.Page,
			PageSize: params.PageSize,
			Total:    total,
		},
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

func (h *PharmacyHandler) toPharmacySummary(ps *models.PharmacySummary) *dtos.PharmacySummary {
	if ps == nil {
		return nil
	}
	return &dtos.PharmacySummary{
		ID:             ps.ID.String(),
		Name:           ps.Name,
		LegalName:      ps.LegalName,
		OnlineHours:    ps.OnlineHours,
		ProfilePicture: ps.ProfilePicture,
		DistanceM:      ps.DistanceM,
		CreatedAt:      ps.CreatedAt,
		UpdatedAt:      ps.UpdatedAt,
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
