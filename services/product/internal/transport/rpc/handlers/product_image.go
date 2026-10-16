package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	"github.com/ritchieridanko/apotekly/services/product/internal/usecases"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductImageHandler struct {
	apis.UnimplementedProductImageServiceServer
	piu usecases.ProductImageUsecase
}

func NewProductImageHandler(piu usecases.ProductImageUsecase) *ProductImageHandler {
	return &ProductImageHandler{piu: piu}
}

func (h *ProductImageHandler) UploadProductImages(ctx context.Context, req *apis.UploadProductImagesRequest) (*apis.UploadProductImagesResponse, error) {
	images := make([]models.ProductImageReq, 0, len(req.GetImages()))
	for _, image := range req.GetImages() {
		if image == nil {
			continue
		}
		images = append(images, *h.toProductImageReq(image))
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

func (h *ProductImageHandler) GetAllProductImages(ctx context.Context, req *apis.GetAllProductImagesRequest) (*apis.GetAllProductImagesResponse, error) {
	pis, err := h.piu.GetAllProductImages(ctx, utils.ToUUID(req.GetProductId()))
	if err != nil {
		return nil, err
	}

	productImages := make([]*apis.ProductImage, 0, len(pis))
	for _, pi := range pis {
		productImages = append(productImages, h.toProductImage(&pi))
	}
	return &apis.GetAllProductImagesResponse{ProductImages: productImages}, nil
}

func (h *ProductImageHandler) UpdateProductImage(ctx context.Context, req *apis.UpdateProductImageRequest) (*apis.UpdateProductImageResponse, error) {
	pi, err := h.piu.UpdateProductImage(
		ctx,
		&models.UpdateProductImageReq{
			ProductImageID: utils.ToUUID(req.GetProductImageId()),
			ProductID:      utils.ToUUID(req.GetProductId()),

			Name:    req.Name,
			AltText: req.AltText,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.UpdateProductImageResponse{ProductImage: h.toProductImage(pi)}, nil
}

func (h *ProductImageHandler) DeleteProductImage(ctx context.Context, req *apis.DeleteProductImageRequest) (*emptypb.Empty, error) {
	err := h.piu.DeleteProductImage(
		ctx,
		&models.DeleteProductImageReq{
			ProductImageID: utils.ToUUID(req.GetProductImageId()),
			ProductID:      utils.ToUUID(req.GetProductId()),
		},
	)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *ProductImageHandler) ReorderProductImage(ctx context.Context, req *apis.ReorderProductImageRequest) (*apis.ReorderProductImageResponse, error) {
	pi, err := h.piu.ReorderProductImage(
		ctx,
		&models.ReorderProductImageReq{
			ProductImageID: utils.ToUUID(req.GetProductImageId()),
			ProductID:      utils.ToUUID(req.GetProductId()),

			NewSortOrder: int16(req.GetNewSortOrder()),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.ReorderProductImageResponse{ProductImage: h.toProductImage(pi)}, nil
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

func (h *ProductImageHandler) toProductImageReq(req *apis.ProductImageRequest) *models.ProductImageReq {
	if req == nil {
		return nil
	}
	return &models.ProductImageReq{
		Name:    req.GetName(),
		URL:     req.GetUrl(),
		AltText: req.AltText,
	}
}
