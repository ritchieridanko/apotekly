package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	"github.com/ritchieridanko/apotekly/services/product/internal/usecases"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
)

type ProductImageHandler struct {
	apis.UnimplementedProductImageServiceServer
	piu usecases.ProductImageUsecase
}

func NewProductImageHandler(piu usecases.ProductImageUsecase) *ProductImageHandler {
	return &ProductImageHandler{piu: piu}
}

func (h *ProductImageHandler) UploadProductImages(ctx context.Context, req *apis.UploadProductImagesRequest) (*apis.UploadProductImagesResponse, error) {
	images := make([]models.UploadProductImageReq, 0, len(req.GetImages()))
	for _, image := range req.GetImages() {
		if image == nil {
			continue
		}
		images = append(images, *h.toUploadProductImageReq(image))
	}

	pis, err := h.piu.UploadProductImages(
		ctx,
		&models.UploadProductImagesReq{
			ProductID: utils.ToUUID(req.GetProductId()),
			Images:    images,
		},
	)
	if err != nil {
		return nil, err
	}

	productImages := make([]*apis.ProductImage, 0, len(pis))
	for _, pi := range pis {
		productImages = append(productImages, h.toProductImage(&pi))
	}
	return &apis.UploadProductImagesResponse{ProductImages: productImages}, nil
}

func (h *ProductImageHandler) toProductImage(pi *models.ProductImage) *apis.ProductImage {
	if pi == nil {
		return nil
	}
	return &apis.ProductImage{
		Id:        pi.ID.String(),
		ProductId: pi.ProductID.String(),
		Name:      pi.Name,
		Image:     pi.Image,
		AltText:   pi.AltText,
		SortOrder: int32(pi.SortOrder),
	}
}

func (h *ProductImageHandler) toUploadProductImageReq(req *apis.UploadProductImageRequest) *models.UploadProductImageReq {
	if req == nil {
		return nil
	}
	return &models.UploadProductImageReq{
		Name:    req.GetName(),
		URL:     req.GetUrl(),
		AltText: req.AltText,
	}
}
