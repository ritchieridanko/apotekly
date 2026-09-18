package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/clients"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/dtos"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type PharmacyHandler struct {
	phc clients.PharmacyClient
}

func NewPharmacyHandler(phc clients.PharmacyClient) *PharmacyHandler {
	return &PharmacyHandler{phc: phc}
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
