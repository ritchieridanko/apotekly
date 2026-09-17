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
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateProductResponse{Product: h.toProduct(p)}, nil
}

func (h *ProductHandler) toProduct(p *models.Product) *apis.Product {
	if p == nil {
		return nil
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
		CreatedAt:      utils.ToTimestamp(&p.CreatedAt),
		UpdatedAt:      utils.ToTimestamp(&p.UpdatedAt),
	}
}
