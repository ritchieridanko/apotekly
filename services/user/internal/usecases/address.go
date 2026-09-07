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
	CreateAddress(ctx context.Context, req *models.CreateAddressReq) (a *models.Address, err *ce.Error)
	GetAllAddresses(ctx context.Context, req *models.GetAllAddressesReq) (as []models.Address, total int64, err *ce.Error)
	UpdateAddress(ctx context.Context, req *models.UpdateAddressReq) (a *models.Address, err *ce.Error)
	SetPrimaryAddress(ctx context.Context, addressID uint64) (a *models.Address, err *ce.Error)
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

func (u *addressUsecase) CreateAddress(ctx context.Context, req *models.CreateAddressReq) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "address.usecase.CreateAddress")
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
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.AddrRecipient(recipient); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if ok, why := u.validator.Phone(phone); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if notes != nil {
		if ok, why := u.validator.AddrNotes(*notes); !ok {
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

	var a *models.Address
	err := u.transactor.WithTx(ctx, func(ctx context.Context) (err *ce.Error) {
		// Primary Address Unsetting (if any)
		if req.IsPrimary {
			_, err = u.ar.UnsetPrimary(ctx, authCtx.AuthID)
			if err != nil {
				return err
			}
		}

		// Address Creation
		a, err = u.ar.Create(
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
		return err
	})
	if err != nil {
		return nil, err.Append(authIDField)
	}

	return a, nil
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
		&models.GetAllAddresses{
			AuthID: authCtx.AuthID,
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

func (u *addressUsecase) UpdateAddress(ctx context.Context, req *models.UpdateAddressReq) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "address.usecase.UpdateAddress")
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
	label := utils.ToLowerPtr(utils.TrimSpacePtr(req.Label))
	recipient := utils.TrimSpacePtr(req.Recipient)
	phone := utils.TrimSpacePtr(req.Phone)
	notes := utils.TrimSpacePtr(req.Notes)
	country := utils.ToLowerPtr(utils.TrimSpacePtr(req.Country))
	subdivision1 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision1))
	subdivision2 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision2))
	subdivision3 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision3))
	subdivision4 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision4))
	street := utils.TrimSpacePtr(req.Street)
	postalCode := utils.ToLowerPtr(utils.TrimSpacePtr(req.PostalCode))

	// Data Validation
	if label != nil {
		if ok, why := u.validator.AddrLabel(*label); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if recipient != nil {
		if ok, why := u.validator.AddrRecipient(*recipient); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if phone != nil {
		if ok, why := u.validator.Phone(*phone); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if notes != nil {
		if ok, why := u.validator.AddrNotes(*notes); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if country != nil {
		if ok, why := u.validator.Country(*country); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
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
	if street != nil {
		if ok, why := u.validator.AddrStreet(*street); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if postalCode != nil {
		if ok, why := u.validator.PostalCode(*postalCode); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if (req.Latitude != nil && req.Longitude == nil) || (req.Latitude == nil && req.Longitude != nil) {
		return nil, ce.NewError(
			ce.CodeInvalidPayload,
			"Both latitude and longitude are required",
			nil,
			authIDField,
		)
	}
	if req.Latitude != nil {
		if ok, why := u.validator.Latitude(*req.Latitude); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if req.Longitude != nil {
		if ok, why := u.validator.Longitude(*req.Longitude); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}

	// Address Update
	a, err := u.ar.Update(
		ctx,
		&models.UpdateAddressP{
			AuthID:    authCtx.AuthID,
			AddressID: req.AddressID,
		},
		&models.UpdateAddressD{
			Label:        label,
			Recipient:    recipient,
			Phone:        phone,
			Notes:        notes,
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
		return nil, err.Append(authIDField)
	}

	return a, nil
}

func (u *addressUsecase) SetPrimaryAddress(ctx context.Context, addressID uint64) (*models.Address, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "address.usecase.SetPrimaryAddress")
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

	var a *models.Address
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Primary Address Unsetting (if any)
		// NOTE: Idempotent scenario (oldPrimary.ID == addressID) does not fail SetPrimaryAddress usecase
		//       But would still return error to cancel the transaction
		oldPrimary, err := u.ar.UnsetPrimary(ctx, authCtx.AuthID)
		if err != nil {
			return err
		}
		if oldPrimary != nil && oldPrimary.ID == addressID {
			a = oldPrimary
			return ce.NewError("idempotent", "", nil)
		}

		// New Primary Address Setting
		a, err = u.ar.SetPrimary(
			ctx,
			&models.SetPrimaryAddress{
				AuthID:    authCtx.AuthID,
				AddressID: addressID,
			},
		)
		return err
	})
	if err != nil && err.Code() != "idempotent" {
		return nil, err.Append(authIDField)
	}

	return a, nil
}
