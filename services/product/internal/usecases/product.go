package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	"github.com/ritchieridanko/apotekly/services/product/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"go.opentelemetry.io/otel"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *models.CreateProductReq) (p *models.Product, err *ce.Error)
}

type productUsecase struct {
	appName   string
	pr        repositories.ProductRepository
	validator *validator.Validator
	logger    *logger.Logger
}

func NewProductUsecase(
	appName string,
	pr repositories.ProductRepository,
	v *validator.Validator,
	l *logger.Logger,
) ProductUsecase {
	return &productUsecase{
		appName:   appName,
		pr:        pr,
		validator: v,
		logger:    l,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, req *models.CreateProductReq) (*models.Product, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "product.usecase.CreateProduct")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}
	if authCtx.PharmacyID == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("pharmacy_id missing from auth context"),
		)
	}

	authFields := []logger.Field{
		logger.NewField("auth_id", authCtx.AuthID),
		logger.NewField("pharmacy_id", authCtx.PharmacyID.String()),
	}

	// Data Normalization
	brandName := strings.TrimSpace(req.BrandName)
	genericName := strings.ToLower(strings.TrimSpace(req.GenericName))
	description := utils.TrimSpacePtr(req.Description)
	dosageForm := strings.ToLower(strings.TrimSpace(req.DosageForm))
	strength := utils.TrimSpacePtr(req.Strength)
	packUnit := strings.ToLower(strings.TrimSpace(req.PackUnit))
	heightCM := utils.Round(req.HeightCM, 2)
	lengthCM := utils.Round(req.LengthCM, 2)
	widthCM := utils.Round(req.WidthCM, 2)
	weightG := utils.Round(req.WeightG, 2)
	price := strings.TrimSpace(req.Price)
	currency := strings.ToLower(strings.TrimSpace(req.Currency))
	manufacturedBy := strings.TrimSpace(req.ManufacturedBy)
	manufacturedIn := strings.ToLower(strings.TrimSpace(req.ManufacturedIn))
	regAuthority := strings.TrimSpace(req.RegAuthority)
	regIdentifier := strings.TrimSpace(req.RegIdentifier)

	// Data Validation
	if ok, why := u.validator.Name(brandName, "Brand name"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Name(genericName, "Generic name"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if description != nil {
		if ok, why := u.validator.Description(*description); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if ok, why := u.validator.DosageForm(dosageForm); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if strength != nil {
		if ok, why := u.validator.DosageStrength(*strength); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if ok, why := u.validator.PackUnit(packUnit); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.PackSize(req.PackSize); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Measurement(heightCM, "Height"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Measurement(lengthCM, "Length"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Measurement(widthCM, "Width"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Measurement(weightG, "Weight"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Currency(currency); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Price(price, currency); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Quantity(req.Quantity); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Name(manufacturedBy, "Manufacturer name"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Country(manufacturedIn, "Country of manufacture"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Name(regAuthority, "Regulatory authority"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.RegIdentifier(regIdentifier); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}

	// Price Parsing
	productPrice, err := utils.ParseMoney(price, currency)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeMoneyParsingFailed,
			ce.MsgInternalServer,
			err,
			authFields...,
		)
	}

	// Product Creation
	p, createErr := u.pr.Create(
		ctx,
		&models.CreateProduct{
			ID:             utils.MustGenerateUUIDv7(),
			PharmacyID:     *authCtx.PharmacyID,
			BrandName:      brandName,
			GenericName:    genericName,
			Description:    description,
			RequiresRX:     req.RequiresRX,
			DosageForm:     dosageForm,
			Strength:       strength,
			PackUnit:       packUnit,
			PackSize:       req.PackSize,
			HeightCM:       heightCM,
			LengthCM:       lengthCM,
			WidthCM:        widthCM,
			WeightG:        weightG,
			Price:          productPrice,
			Currency:       currency,
			Quantity:       req.Quantity,
			IsActive:       req.IsActive,
			ManufacturedBy: manufacturedBy,
			ManufacturedIn: manufacturedIn,
			RegAuthority:   regAuthority,
			RegIdentifier:  regIdentifier,
		},
	)
	if createErr != nil {
		return nil, createErr.Append(authFields...)
	}

	return p, nil
}
