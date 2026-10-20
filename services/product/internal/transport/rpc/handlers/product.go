package handlers

import (
	"context"
	"strings"

	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	"github.com/ritchieridanko/apotekly/services/product/internal/usecases"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
)

type ProductHandler struct {
	apis.UnimplementedProductServiceServer
	pu usecases.ProductUsecase
}

func NewProductHandler(pu usecases.ProductUsecase) *ProductHandler {
	return &ProductHandler{pu: pu}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *apis.CreateProductRequest) (*apis.CreateProductResponse, error) {
	images := make([]models.ProductImageReq, 0, len(req.GetImages()))
	for _, image := range req.GetImages() {
		if image == nil {
			continue
		}
		images = append(images, *h.toProductImageReq(image))
	}

	p, err := h.pu.CreateProduct(
		ctx,
		&models.CreateProductReq{
			BrandName:      req.GetBrandName(),
			GenericName:    req.GetGenericName(),
			Description:    req.Description,
			RequiresRX:     req.GetRequiresRx(),
			DosageForm:     req.GetDosageForm(),
			Strength:       req.Strength,
			PackUnit:       req.GetPackUnit(),
			PackSize:       int(req.GetPackSize()),
			HeightCM:       req.GetHeightCm(),
			LengthCM:       req.GetLengthCm(),
			WidthCM:        req.GetWidthCm(),
			WeightG:        req.GetWeightG(),
			Price:          req.GetPrice(),
			Currency:       req.GetCurrency(),
			Quantity:       int(req.GetQuantity()),
			IsActive:       req.GetIsActive(),
			ManufacturedBy: req.GetManufacturedBy(),
			ManufacturedIn: req.GetManufacturedIn(),
			RegAuthority:   req.GetRegAuthority(),
			RegIdentifier:  req.GetRegIdentifier(),
			Images:         images,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateProductResponse{Product: h.toProduct(p)}, nil
}

func (h *ProductHandler) GetProductByID(ctx context.Context, req *apis.GetProductByIDRequest) (*apis.GetProductByIDResponse, error) {
	p, err := h.pu.GetProductByID(ctx, utils.ToUUID(req.GetProductId()))
	if err != nil {
		return nil, err
	}
	return &apis.GetProductByIDResponse{Product: h.toProduct(p)}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *apis.UpdateProductRequest) (*apis.UpdateProductResponse, error) {
	p, err := h.pu.UpdateProduct(
		ctx,
		&models.UpdateProductReq{
			ProductID: utils.ToUUID(req.GetProductId()),

			BrandName:      req.BrandName,
			GenericName:    req.GenericName,
			Description:    req.Description,
			RequiresRX:     req.RequiresRx,
			DosageForm:     req.DosageForm,
			Strength:       req.Strength,
			PackUnit:       req.PackUnit,
			PackSize:       utils.ToIntPtr(req.PackSize),
			HeightCM:       req.HeightCm,
			LengthCM:       req.LengthCm,
			WidthCM:        req.WidthCm,
			WeightG:        req.WeightG,
			Price:          req.Price,
			Currency:       req.Currency,
			Quantity:       utils.ToIntPtr(req.Quantity),
			ManufacturedBy: req.ManufacturedBy,
			ManufacturedIn: req.ManufacturedIn,
			RegAuthority:   req.RegAuthority,
			RegIdentifier:  req.RegIdentifier,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.UpdateProductResponse{Product: h.toProduct(p)}, nil
}

func (h *ProductHandler) toProduct(p *models.Product) *apis.Product {
	if p == nil {
		return nil
	}

	images := make([]*apis.ProductImage, 0, len(p.Images))
	for _, image := range p.Images {
		images = append(images, h.toProductImage(&image))
	}

	return &apis.Product{
		Id:             p.ID.String(),
		PharmacyId:     p.PharmacyID.String(),
		BrandName:      p.BrandName,
		GenericName:    utils.ToTitlecase(p.GenericName),
		Description:    p.Description,
		RequiresRx:     p.RequiresRX,
		DosageForm:     utils.ToTitlecase(p.DosageForm),
		Strength:       p.Strength,
		PackUnit:       utils.ToTitlecase(p.PackUnit),
		PackSize:       int32(p.PackSize),
		HeightCm:       p.HeightCM,
		LengthCm:       p.LengthCM,
		WidthCm:        p.WidthCM,
		WeightG:        p.WeightG,
		Price:          utils.StringifyMoney(p.Price, p.Currency),
		Currency:       strings.ToUpper(p.Currency),
		Quantity:       int32(p.Quantity),
		IsActive:       p.IsActive,
		ManufacturedBy: p.ManufacturedBy,
		ManufacturedIn: utils.ToTitlecase(p.ManufacturedIn),
		RegAuthority:   p.RegAuthority,
		RegIdentifier:  p.RegIdentifier,
		Images:         images,
		CreatedAt:      utils.ToTimestamp(&p.CreatedAt),
		UpdatedAt:      utils.ToTimestamp(&p.UpdatedAt),
	}
}

func (h *ProductHandler) toProductImage(pi *models.ProductImage) *apis.ProductImage {
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

func (h *ProductHandler) toProductImageReq(req *apis.ProductImageRequest) *models.ProductImageReq {
	if req == nil {
		return nil
	}
	return &models.ProductImageReq{
		Name:    req.GetName(),
		URL:     req.GetUrl(),
		AltText: req.AltText,
	}
}
