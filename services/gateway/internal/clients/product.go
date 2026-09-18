package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var productServiceField logger.Field = logger.NewField("service", "product")

type ProductClient interface {
	CreateProduct(ctx context.Context, req *models.CreateProductReq) (p *models.Product, err *ce.Error)
}

type productClient struct {
	client apis.ProductServiceClient
}

func NewProductClient(c apis.ProductServiceClient) ProductClient {
	return &productClient{client: c}
}

func (c *productClient) CreateProduct(ctx context.Context, req *models.CreateProductReq) (*models.Product, *ce.Error) {
	resp, err := c.client.CreateProduct(
		ctx,
		&apis.CreateProductRequest{
			BrandName:      req.BrandName,
			GenericName:    req.GenericName,
			Description:    req.Description,
			RequiresRx:     req.RequiresRX,
			DosageForm:     req.DosageForm,
			Strength:       req.Strength,
			PackUnit:       req.PackUnit,
			PackSize:       int32(req.PackSize),
			HeightCm:       req.HeightCM,
			LengthCm:       req.LengthCM,
			WidthCm:        req.WidthCM,
			WeightG:        req.WeightG,
			Price:          req.Price,
			Currency:       req.Currency,
			Quantity:       int32(req.Quantity),
			IsActive:       req.IsActive,
			ManufacturedBy: req.ManufacturedBy,
			ManufacturedIn: req.ManufacturedIn,
			RegAuthority:   req.RegAuthority,
			RegIdentifier:  req.RegIdentifier,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			productServiceField,
		)
	}
	return c.toProduct(resp.GetProduct()), nil
}

func (c *productClient) toProduct(p *apis.Product) *models.Product {
	if p == nil {
		return nil
	}
	return &models.Product{
		ID:             utils.ToUUID(p.GetId()),
		PharmacyID:     utils.ToUUID(p.GetPharmacyId()),
		BrandName:      p.GetBrandName(),
		GenericName:    p.GetGenericName(),
		Description:    p.Description,
		RequiresRX:     p.GetRequiresRx(),
		DosageForm:     p.GetDosageForm(),
		Strength:       p.Strength,
		PackUnit:       p.GetPackUnit(),
		PackSize:       int(p.GetPackSize()),
		HeightCM:       p.GetHeightCm(),
		LengthCM:       p.GetLengthCm(),
		WidthCM:        p.GetWidthCm(),
		WeightG:        p.GetWeightG(),
		Price:          p.GetPrice(),
		Currency:       p.GetCurrency(),
		Quantity:       int(p.GetQuantity()),
		IsActive:       p.GetIsActive(),
		ManufacturedBy: p.GetManufacturedBy(),
		ManufacturedIn: p.GetManufacturedIn(),
		RegAuthority:   p.GetRegAuthority(),
		RegIdentifier:  p.GetRegIdentifier(),
		CreatedAt:      utils.ToTime(p.GetCreatedAt()),
		UpdatedAt:      utils.ToTime(p.GetUpdatedAt()),
	}
}
