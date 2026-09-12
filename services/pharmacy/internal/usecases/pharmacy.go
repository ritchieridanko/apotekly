package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/models"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"go.opentelemetry.io/otel"
)

type PharmacyUsecase interface {
	CreatePharmacy(ctx context.Context, req *models.CreatePharmacyReq) (p *models.Pharmacy, err *ce.Error)
	GetMe(ctx context.Context) (p *models.Pharmacy, err *ce.Error)
}

type pharmacyUsecase struct {
	appName   string
	pr        repositories.PharmacyRepository
	validator *validator.Validator
	logger    *logger.Logger
}

func NewPharmacyUsecase(
	appName string,
	pr repositories.PharmacyRepository,
	v *validator.Validator,
	l *logger.Logger,
) PharmacyUsecase {
	return &pharmacyUsecase{
		appName:   appName,
		pr:        pr,
		validator: v,
		logger:    l,
	}
}

func (u *pharmacyUsecase) CreatePharmacy(ctx context.Context, req *models.CreatePharmacyReq) (*models.Pharmacy, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "pharmacy.usecase.CreatePharmacy")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Data Normalization
	name := strings.TrimSpace(req.Name)
	legalName := utils.TrimSpacePtr(req.LegalName)
	description := utils.TrimSpacePtr(req.Description)
	country := strings.ToLower(strings.TrimSpace(req.Country))
	subdivision1 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision1))
	subdivision2 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision2))
	subdivision3 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision3))
	subdivision4 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision4))
	street := strings.TrimSpace(req.Street)
	postalCode := strings.ToLower(strings.TrimSpace(req.PostalCode))
	email := utils.ToLowerPtr(utils.TrimSpacePtr(req.Email))
	phone := utils.TrimSpacePtr(req.Phone)
	website := utils.TrimSpacePtr(req.Website)
	whatsapp := utils.TrimSpacePtr(req.Whatsapp)

	// Data Validation
	if ok, why := u.validator.Name(name); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if legalName != nil {
		if ok, why := u.validator.Name(*legalName); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if description != nil {
		if ok, why := u.validator.Description(*description); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if req.OnlineHours != nil {
		if ok, why := u.validator.OnlineHours(*req.OnlineHours); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if ok, why := u.validator.Country(country); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if subdivision1 != nil {
		if ok, why := u.validator.AddrSubdivision(*subdivision1); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if subdivision2 != nil {
		if ok, why := u.validator.AddrSubdivision(*subdivision2); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if subdivision3 != nil {
		if ok, why := u.validator.AddrSubdivision(*subdivision3); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if subdivision4 != nil {
		if ok, why := u.validator.AddrSubdivision(*subdivision4); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if ok, why := u.validator.AddrStreet(street); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.PostalCode(postalCode); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.Latitude(req.Latitude); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.Longitude(req.Longitude); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if email != nil {
		if ok, why := u.validator.Email(*email); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if phone != nil {
		if ok, why := u.validator.Phone(*phone); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if website != nil {
		if ok, why := u.validator.URL(*website); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if whatsapp != nil {
		if ok, why := u.validator.Phone(*whatsapp); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}

	// Pharmacy Creation
	p, err := u.pr.Create(
		ctx,
		&models.CreatePharmacy{
			ID:           utils.MustGenerateUUIDv7(),
			AuthID:       authCtx.AuthID,
			Name:         name,
			LegalName:    legalName,
			Description:  description,
			OnlineHours:  req.OnlineHours,
			Country:      country,
			Subdivision1: subdivision1,
			Subdivision2: subdivision2,
			Subdivision3: subdivision3,
			Subdivision4: subdivision4,
			Street:       street,
			PostalCode:   postalCode,
			Latitude:     req.Latitude,
			Longitude:    req.Longitude,
			Email:        email,
			Phone:        phone,
			Website:      website,
			Whatsapp:     whatsapp,
		},
	)
	if err != nil {
		return nil, err.Append(authIDField)
	}

	return p, nil
}

func (u *pharmacyUsecase) GetMe(ctx context.Context) (*models.Pharmacy, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "pharmacy.usecase.GetMe")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Pharmacy Fetching
	p, err := u.pr.GetByAuthID(ctx, authCtx.AuthID)
	if err != nil {
		return nil, err.Append(authIDField)
	}
	return p, nil
}
