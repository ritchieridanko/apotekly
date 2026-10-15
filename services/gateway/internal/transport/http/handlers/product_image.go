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

type ProductImageHandler struct {
	pic clients.ProductImageClient
}

func NewProductImageHandler(pic clients.ProductImageClient) *ProductImageHandler {
	return &ProductImageHandler{pic: pic}
}

func (h *ProductImageHandler) UploadProductImages(ctx *gin.Context) {
	var payload dtos.UploadProductImagesRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	images := make([]models.UploadProductImageReq, 0, len(payload.Images))
	for _, image := range payload.Images {
		images = append(images, *h.toUploadProductImageReq(&image))
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

	pis, err := h.pic.UploadProductImages(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID, strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole, authCtx.Role,
			constants.MDKeyPharmacyID, authCtx.PharmacyID.String(),
		),
		&models.UploadProductImagesReq{
			ProductID: ctx.Param(constants.URLParamProductID),

			Images: images,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	productImages := make([]dtos.ProductImage, 0, len(pis))
	for _, pi := range pis {
		productImages = append(productImages, *h.toProductImage(&pi))
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusCreated,
		"Product images uploaded successfully",
		dtos.UploadProductImagesResponse{ProductImages: productImages},
		nil,
	)
}

func (h *ProductImageHandler) GetAllProductImages(ctx *gin.Context) {
	pis, err := h.pic.GetAllProductImages(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
		),
		ctx.Param(constants.URLParamProductID),
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	productImages := make([]dtos.ProductImage, 0, len(pis))
	for _, pi := range pis {
		productImages = append(productImages, *h.toProductImage(&pi))
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Product images retrieved successfully",
		dtos.GetAllProductImagesResponse{ProductImages: productImages},
		nil,
	)
}

func (h *ProductImageHandler) UpdateProductImage(ctx *gin.Context) {
	var payload dtos.UpdateProductImageRequest
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

	pi, err := h.pic.UpdateProductImage(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID, strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole, authCtx.Role,
			constants.MDKeyPharmacyID, authCtx.PharmacyID.String(),
		),
		&models.UpdateProductImageReq{
			ProductImageID: ctx.Param(constants.URLParamProductImageID),
			ProductID:      ctx.Param(constants.URLParamProductID),

			Name:    payload.Name,
			AltText: payload.AltText,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Product image updated successfully",
		dtos.UpdateProductImageResponse{ProductImage: h.toProductImage(pi)},
		nil,
	)
}

func (h *ProductImageHandler) DeleteProductImage(ctx *gin.Context) {
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

	err := h.pic.DeleteProductImage(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID, strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole, authCtx.Role,
			constants.MDKeyPharmacyID, authCtx.PharmacyID.String(),
		),
		&models.DeleteProductImageReq{
			ProductImageID: ctx.Param(constants.URLParamProductImageID),
			ProductID:      ctx.Param(constants.URLParamProductID),
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse[any](ctx, http.StatusNoContent, "", nil, nil)
}

func (h *ProductImageHandler) ReorderProductImage(ctx *gin.Context) {
	var payload dtos.ReorderProductImageRequest
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

	pi, err := h.pic.ReorderProductImage(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID, strconv.FormatUint(authCtx.AuthID, 10),
			constants.MDKeyRole, authCtx.Role,
			constants.MDKeyPharmacyID, authCtx.PharmacyID.String(),
		),
		&models.ReorderProductImageReq{
			ProductImageID: ctx.Param(constants.URLParamProductImageID),
			ProductID:      ctx.Param(constants.URLParamProductID),

			NewSortOrder: payload.NewSortOrder,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Product image reordered successfully",
		dtos.ReorderProductImageResponse{ProductImage: h.toProductImage(pi)},
		nil,
	)
}

func (h *ProductImageHandler) toProductImage(pi *models.ProductImage) *dtos.ProductImage {
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

func (h *ProductImageHandler) toUploadProductImageReq(req *dtos.UploadProductImageRequest) *models.UploadProductImageReq {
	if req == nil {
		return nil
	}
	return &models.UploadProductImageReq{
		Name:    req.Name,
		URL:     req.URL,
		AltText: req.AltText,
	}
}
