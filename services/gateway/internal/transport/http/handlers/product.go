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

	images := make([]models.ProductImageReq, 0, len(payload.Images))
	for _, image := range payload.Images {
		images = append(images, *h.toProductImageReq(&image))
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
			constants.MDKeyAuthID, strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole, authCtx.Role,
			constants.MDKeyPharmacyID, authCtx.PharmacyID.String(),
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
			Images:         images,
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

func (h *ProductHandler) GetProductByID(ctx *gin.Context) {
	p, err := h.prc.GetProductByID(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
		),
		ctx.Param(constants.URLParamProductID),
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Product retrieved successfully",
		dtos.GetProductByIDResponse{Product: h.toProduct(p)},
		nil,
	)
}

func (h *ProductHandler) toProduct(p *models.Product) *dtos.Product {
	if p == nil {
		return nil
	}

	images := make([]dtos.ProductImage, 0, len(p.Images))
	for _, image := range p.Images {
		images = append(images, *h.toProductImage(&image))
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
		Images:         images,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func (h *ProductHandler) toProductImage(pi *models.ProductImage) *dtos.ProductImage {
	if pi == nil {
		return nil
	}
	return &dtos.ProductImage{
		ID:        pi.ID.String(),
		ProductID: pi.ProductID.String(),
		Name:      pi.Name,
		Image:     pi.Image,
		AltText:   pi.AltText,
		SortOrder: pi.SortOrder,
	}
}

func (h *ProductHandler) toProductImageReq(req *dtos.ProductImageRequest) *models.ProductImageReq {
	if req == nil {
		return nil
	}
	return &models.ProductImageReq{
		Name:    req.Name,
		URL:     req.URL,
		AltText: req.AltText,
	}
}
