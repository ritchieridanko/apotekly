package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
	"github.com/ritchieridanko/apotekly/services/user/internal/repositories"
	"go.opentelemetry.io/otel"
)

const (
	defaultPageSize int = 10
	maxPageSize     int = 100
)

type AddressUsecase interface {
	CreateAddress(ctx context.Context, req *models.CreateAddressReq) (a *models.Address, oldPrimary *models.Address, err *ce.Error)
	GetAllAddresses(ctx context.Context, req *models.GetAllAddressesReq) (as []models.Address, total int64, err *ce.Error)
}

type addressUsecase struct {
	appName    string
	ar         repositories.AddressRepository
	transactor *database.Transactor
	validator  *validator.Validator
	logger     *logger.Logger
}

func NewAddressUsecase(
	appName string,
	ar repositories.AddressRepository,
	tx *database.Transactor,
	v *validator.Validator,
	l *logger.Logger,
) AddressUsecase {
	return &addressUsecase{
		appName:    appName,
		ar:         ar,
		transactor: tx,
		validator:  v,
		logger:     l,
	}
}

func (u *addressUsecase) CreateAddress(ctx context.Context, req *models.CreateAddressReq) (*models.Address, *models.Address, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "address.usecase.CreateAddress")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Data Normalization
	label := strings.ToLower(strings.TrimSpace(req.Label))
	recipient := strings.TrimSpace(req.Recipient)
	phone := strings.TrimSpace(req.Phone)
	notes := utils.TrimSpacePtr(req.Notes)
	country := strings.ToLower(strings.TrimSpace(req.Country))
	subdivision1 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision1))
	subdivision2 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision2))
	subdivision3 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision3))
	subdivision4 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision4))
	street := strings.TrimSpace(req.Street)
	postalCode := strings.ToLower(strings.TrimSpace(req.PostalCode))

	// Data Validation
	if ok, why := u.validator.AddrLabel(label); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.AddrRecipient(recipient); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.Phone(phone); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if notes != nil {
		if ok, why := u.validator.AddrNotes(*notes); !ok {
			return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if ok, why := u.validator.Country(country); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if subdivision1 != nil {
		if ok, why := u.validator.AddrSubdivision(*subdivision1); !ok {
			return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if subdivision2 != nil {
		if ok, why := u.validator.AddrSubdivision(*subdivision2); !ok {
			return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if subdivision3 != nil {
		if ok, why := u.validator.AddrSubdivision(*subdivision3); !ok {
			return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if subdivision4 != nil {
		if ok, why := u.validator.AddrSubdivision(*subdivision4); !ok {
			return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if ok, why := u.validator.AddrStreet(street); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.PostalCode(postalCode); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.Latitude(req.Latitude); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.Longitude(req.Longitude); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}

	var a, oldPrimaryAddress *models.Address
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Primary Address Unsetting (if any)
		if req.IsPrimary {
			address, err := u.ar.UnsetPrimary(ctx, authCtx.AuthID)
			if err != nil {
				return err
			}
			if address != nil {
				oldPrimaryAddress = address
			}
		}

		// Address Creation
		address, err := u.ar.Create(
			ctx,
			&models.CreateAddress{
				AuthID:       authCtx.AuthID,
				Label:        label,
				Recipient:    recipient,
				Phone:        phone,
				Notes:        notes,
				IsPrimary:    req.IsPrimary,
				Country:      country,
				Subdivision1: subdivision1,
				Subdivision2: subdivision2,
				Subdivision3: subdivision3,
				Subdivision4: subdivision4,
				Street:       street,
				PostalCode:   postalCode,
				Latitude:     req.Latitude,
				Longitude:    req.Longitude,
			},
		)
		if err != nil {
			return err
		}

		a = address
		return nil
	})
	if err != nil {
		return nil, nil, err.Append(authIDField)
	}

	return a, oldPrimaryAddress, nil
}

func (u *addressUsecase) GetAllAddresses(ctx context.Context, req *models.GetAllAddressesReq) ([]models.Address, int64, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "address.usecase.GetAllAddresses")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, 0, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Data Validation
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	// All Addresses Fetching
	as, total, err := u.ar.GetAll(
		ctx,
		authCtx.AuthID,
		&models.GetAllAddresses{
			OffsetPagination: utils.OffsetPagination{
				Page:     page,
				PageSize: pageSize,
			},
		},
	)
	if err != nil {
		return nil, 0, err.Append(authIDField)
	}

	return as, total, nil
}
