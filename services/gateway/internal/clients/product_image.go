package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var productImageServiceField logger.Field = logger.NewField("service", "product.image")

type ProductImageClient interface {
	UploadProductImages(ctx context.Context, req *models.UploadProductImagesReq) (pis []models.ProductImage, err *ce.Error)
	GetAllProductImages(ctx context.Context, productID string) (pis []models.ProductImage, err *ce.Error)
	UpdateProductImage(ctx context.Context, req *models.UpdateProductImageReq) (pi *models.ProductImage, err *ce.Error)
	DeleteProductImage(ctx context.Context, req *models.DeleteProductImageReq) (err *ce.Error)
	ReorderProductImage(ctx context.Context, req *models.ReorderProductImageReq) (pi *models.ProductImage, err *ce.Error)
}

type productImageClient struct {
	client apis.ProductImageServiceClient
}

func NewProductImageClient(c apis.ProductImageServiceClient) ProductImageClient {
	return &productImageClient{client: c}
}

func (c *productImageClient) UploadProductImages(ctx context.Context, req *models.UploadProductImagesReq) ([]models.ProductImage, *ce.Error) {
	images := make([]*apis.ProductImageRequest, 0, len(req.Images))
	for _, image := range req.Images {
		images = append(images, c.toProductImageRequest(&image))
	}

	resp, err := c.client.UploadProductImages(
		ctx,
		&apis.UploadProductImagesRequest{
			ProductId: req.ProductID,
			Images:    images,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			productImageServiceField,
		)
	}

	pis := make([]models.ProductImage, 0, len(resp.GetProductImages()))
	for _, pi := range resp.GetProductImages() {
		if pi == nil {
			continue
		}
		pis = append(pis, *c.toProductImage(pi))
	}

	return pis, nil
}

func (c *productImageClient) GetAllProductImages(ctx context.Context, productID string) ([]models.ProductImage, *ce.Error) {
	resp, err := c.client.GetAllProductImages(
		ctx,
		&apis.GetAllProductImagesRequest{
			ProductId: productID,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			productImageServiceField,
		)
	}

	pis := make([]models.ProductImage, 0, len(resp.GetProductImages()))
	for _, pi := range resp.GetProductImages() {
		if pi == nil {
			continue
		}
		pis = append(pis, *c.toProductImage(pi))
	}

	return pis, nil
}

func (c *productImageClient) UpdateProductImage(ctx context.Context, req *models.UpdateProductImageReq) (*models.ProductImage, *ce.Error) {
	resp, err := c.client.UpdateProductImage(
		ctx,
		&apis.UpdateProductImageRequest{
			ProductImageId: req.ProductImageID,
			ProductId:      req.ProductID,

			Name:    req.Name,
			AltText: req.AltText,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			productImageServiceField,
		)
	}
	return c.toProductImage(resp.GetProductImage()), nil
}

func (c *productImageClient) DeleteProductImage(ctx context.Context, req *models.DeleteProductImageReq) *ce.Error {
	_, err := c.client.DeleteProductImage(
		ctx,
		&apis.DeleteProductImageRequest{
			ProductImageId: req.ProductImageID,
			ProductId:      req.ProductID,
		},
	)
	if err != nil {
		return ce.ToError(
			err,
		).Append(
			productImageServiceField,
		)
	}
	return nil
}

func (c *productImageClient) ReorderProductImage(ctx context.Context, req *models.ReorderProductImageReq) (*models.ProductImage, *ce.Error) {
	resp, err := c.client.ReorderProductImage(
		ctx,
		&apis.ReorderProductImageRequest{
			ProductImageId: req.ProductImageID,
			ProductId:      req.ProductID,

			NewSortOrder: int32(req.NewSortOrder),
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			productImageServiceField,
		)
	}
	return c.toProductImage(resp.GetProductImage()), nil
}

func (c *productImageClient) toProductImage(pi *apis.ProductImage) *models.ProductImage {
	if pi == nil {
		return nil
	}
	return &models.ProductImage{
		ID:        utils.ToUUID(pi.GetId()),
		ProductID: utils.ToUUID(pi.GetProductId()),
		Name:      pi.GetName(),
		Image:     pi.GetImage(),
		AltText:   pi.AltText,
		SortOrder: int16(pi.GetSortOrder()),
	}
}

func (c *productImageClient) toProductImageRequest(req *models.ProductImageReq) *apis.ProductImageRequest {
	if req == nil {
		return nil
	}
	return &apis.ProductImageRequest{
		Name:    req.Name,
		Url:     req.URL,
		AltText: req.AltText,
	}
}
