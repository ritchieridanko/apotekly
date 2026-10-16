package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/product/internal/clients"
	"github.com/ritchieridanko/apotekly/services/product/internal/models"
	"github.com/ritchieridanko/apotekly/services/product/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"go.opentelemetry.io/otel"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *models.CreateProductReq) (p *models.Product, err *ce.Error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (p *models.Product, err *ce.Error)
}

type productUsecase struct {
	appName    string
	piu        ProductImageUsecase
	pr         repositories.ProductRepository
	pc         clients.PharmacyClient
	transactor *database.Transactor
	validator  *validator.Validator
	logger     *logger.Logger
}

func NewProductUsecase(
	appName string,
	piu ProductImageUsecase,
	pr repositories.ProductRepository,
	pc clients.PharmacyClient,
	tx *database.Transactor,
	v *validator.Validator,
	l *logger.Logger,
) ProductUsecase {
	return &productUsecase{
		appName:    appName,
		piu:        piu,
		pr:         pr,
		pc:         pc,
		transactor: tx,
		validator:  v,
		logger:     l,
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

	// Pharmacy Status Fetching
	status, err := u.pc.GetStatus(ctx, authCtx.AuthID)
	if err != nil && err.Code() == ce.CodeNotFound {
		return nil, ce.NewError(
			ce.CodeRoleMismatch,
			ce.MsgInternalServer,
			err.Unwrap(),
			err.Append(
				authFields...,
			).Fields()...,
		)
	}
	if err != nil {
		return nil, err.Append(authFields...)
	}
	if status != constants.PharmacyStatusActive {
		return nil, ce.NewError(
			ce.CodePharmacyNotActive,
			"Pharmacy account is "+status,
			nil,
			authFields...,
		)
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
	productPrice, parseErr := utils.ParseMoney(price, currency)
	if parseErr != nil {
		return nil, ce.NewError(
			ce.CodeMoneyParsingFailed,
			ce.MsgInternalServer,
			parseErr,
			authFields...,
		)
	}

	var p *models.Product
	err = u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Product Creation
		p, err = u.pr.Create(
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
		if err != nil {
			return err
		}

		if len(req.Images) > 0 {
			// Product Images Creation
			pis, err := u.piu.CreateProductImages(
				ctx,
				&models.CreateProductImagesReq{
					ProductID: p.ID,
					Images:    req.Images,
				},
			)
			if err != nil {
				return err
			}

			p.Images = pis
		}

		return nil
	})
	if err != nil {
		return nil, err.Append(authFields...)
	}

	return p, nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, productID uuid.UUID) (*models.Product, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "product.usecase.GetProductByID")
	defer span.End()

	// Product Fetching
	p, err := u.pr.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// All Product Images Fetching
	pis, err := u.piu.GetAllProductImages(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	p.Images = pis
	return p, nil
}
