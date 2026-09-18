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

type ProductHandler struct {
	prc clients.ProductClient
}

func NewProductHandler(prc clients.ProductClient) *ProductHandler {
	return &ProductHandler{prc: prc}
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var payload dtos.CreateProductRequest
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

	p, err := h.prc.CreateProduct(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole,
			authCtx.Role,
			constants.MDKeyPharmacyID,
			authCtx.PharmacyID.String(),
		),
		&models.CreateProductReq{
			BrandName:      payload.BrandName,
			GenericName:    payload.GenericName,
			Description:    payload.Description,
			RequiresRX:     payload.RequiresRX,
			DosageForm:     payload.DosageForm,
			Strength:       payload.Strength,
			PackUnit:       payload.PackUnit,
			PackSize:       payload.PackSize,
			HeightCM:       payload.HeightCM,
			LengthCM:       payload.LengthCM,
			WidthCM:        payload.WidthCM,
			WeightG:        payload.WeightG,
			Price:          payload.Price,
			Currency:       payload.Currency,
			Quantity:       payload.Quantity,
			IsActive:       payload.IsActive,
			ManufacturedBy: payload.ManufacturedBy,
			ManufacturedIn: payload.ManufacturedIn,
			RegAuthority:   payload.RegAuthority,
			RegIdentifier:  payload.RegIdentifier,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusCreated,
		"Product created successfully",
		dtos.CreateProductResponse{Product: h.toProduct(p)},
		nil,
	)
}

func (h *ProductHandler) toProduct(p *models.Product) *dtos.Product {
	if p == nil {
		return nil
	}
	return &dtos.Product{
		ID:             p.ID.String(),
		PharmacyID:     p.PharmacyID.String(),
		BrandName:      p.BrandName,
		GenericName:    p.GenericName,
		Description:    p.Description,
		RequiresRX:     p.RequiresRX,
		DosageForm:     p.DosageForm,
		Strength:       p.Strength,
		PackUnit:       p.PackUnit,
		PackSize:       p.PackSize,
		HeightCM:       p.HeightCM,
		LengthCM:       p.LengthCM,
		WidthCM:        p.WidthCM,
		WeightG:        p.WeightG,
		Price:          p.Price,
		Currency:       p.Currency,
		Quantity:       p.Quantity,
		IsActive:       p.IsActive,
		ManufacturedBy: p.ManufacturedBy,
		ManufacturedIn: p.ManufacturedIn,
		RegAuthority:   p.RegAuthority,
		RegIdentifier:  p.RegIdentifier,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
