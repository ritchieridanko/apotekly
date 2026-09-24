package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/clients"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/models"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/storage"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"go.opentelemetry.io/otel"
)

type PharmacyUsecase interface {
	CreatePharmacy(ctx context.Context, req *models.CreatePharmacyReq) (p *models.Pharmacy, err *ce.Error)
	GetMe(ctx context.Context) (p *models.Pharmacy, err *ce.Error)
	GetID(ctx context.Context, authID uint64) (pharmacyID uuid.UUID, err *ce.Error)
	GetPharmacyByID(ctx context.Context, pharmacyID uuid.UUID) (p *models.Pharmacy, err *ce.Error)
	GetAllPharmacies(ctx context.Context, req *models.GetAllPharmaciesReq) (pss []models.PharmacySummary, total int64, err *ce.Error)
	UpdatePharmacy(ctx context.Context, req *models.UpdatePharmacyReq) (p *models.Pharmacy, err *ce.Error)
	UpdateProfilePicture(ctx context.Context, profilePictureURL string) (p *models.Pharmacy, err *ce.Error)
	UpdateProfileBanner(ctx context.Context, profileBannerURL string) (p *models.Pharmacy, err *ce.Error)
}

type pharmacyUsecase struct {
	appName    string
	pr         repositories.PharmacyRepository
	ac         clients.AuthClient
	uac        clients.AddressClient
	transactor *database.Transactor
	storage    *storage.Storage
	validator  *validator.Validator
	logger     *logger.Logger
}

func NewPharmacyUsecase(
	appName string,
	pr repositories.PharmacyRepository,
	ac clients.AuthClient,
	uac clients.AddressClient,
	tx *database.Transactor,
	s *storage.Storage,
	v *validator.Validator,
	l *logger.Logger,
) PharmacyUsecase {
	return &pharmacyUsecase{
		appName:    appName,
		pr:         pr,
		ac:         ac,
		uac:        uac,
		transactor: tx,
		storage:    s,
		validator:  v,
		logger:     l,
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
	if ok, why := u.validator.Name(name, "Name"); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if legalName != nil {
		if ok, why := u.validator.Name(*legalName, "Legal name"); !ok {
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
	if ok, why := u.validator.Country(country, "Country"); !ok {
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
	if ok, why := u.validator.Name(street, "Street name"); !ok {
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

	var p *models.Pharmacy
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Pharmacy Creation
		pharmacy, err := u.pr.Create(
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
			return err
		}

		// New Role Setting
		if err := u.ac.SetRolePharmacy(ctx, authCtx.AuthID); err != nil {
			if err.Code() == ce.CodeNotFound {
				return ce.NewError(
					ce.CodeAuthNotRegistered,
					ce.MsgInvalidCredentials,
					err.Unwrap(),
					err.Fields()...,
				)
			}
			return err
		}

		p = pharmacy
		return nil
	})
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

func (u *pharmacyUsecase) GetID(ctx context.Context, authID uint64) (uuid.UUID, *ce.Error) {
	return u.pr.GetID(ctx, authID)
}

func (u *pharmacyUsecase) GetPharmacyByID(ctx context.Context, pharmacyID uuid.UUID) (*models.Pharmacy, *ce.Error) {
	return u.pr.GetByID(ctx, pharmacyID)
}

func (u *pharmacyUsecase) GetAllPharmacies(ctx context.Context, req *models.GetAllPharmaciesReq) ([]models.PharmacySummary, int64, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "pharmacy.usecase.GetAllPharmacies")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)

	// Data Normalization
	search := utils.TrimSpacePtr(req.Search)

	// Data Validation
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = constants.PageDefaultSizePharmacy
	}
	if pageSize > constants.PageMaxSizePharmacy {
		pageSize = constants.PageMaxSizePharmacy
	}
	if search != nil {
		if ok, why := u.validator.Search(*search); !ok {
			return nil, 0, ce.NewError(ce.CodeInvalidParams, why, nil)
		}
	}
	if req.RadiusM != nil {
		if ok, why := u.validator.Radius(*req.RadiusM); !ok {
			return nil, 0, ce.NewError(ce.CodeInvalidParams, why, nil)
		}
	}
	if req.RequireLocation() {
		if req.Latitude != nil && req.Longitude != nil {
			if ok, why := u.validator.Latitude(*req.Latitude); !ok {
				return nil, 0, ce.NewError(ce.CodeInvalidParams, why, nil)
			}
			if ok, why := u.validator.Longitude(*req.Longitude); !ok {
				return nil, 0, ce.NewError(ce.CodeInvalidParams, why, nil)
			}
		} else {
			// Authentication Status Check
			if authCtx == nil {
				return nil, 0, ce.NewError(
					ce.CodeLocationNotProvided,
					ce.MsgLocationNotProvided,
					nil,
				)
			}

			// Primary Location Fetching
			lat, lon, err := u.uac.GetPrimaryLocation(ctx, authCtx.AuthID)
			if err != nil && err.Code() == ce.CodeAddressNotFound {
				return nil, 0, ce.NewError(
					ce.CodeLocationNotProvided,
					ce.MsgLocationNotProvided,
					err.Unwrap(),
					err.Fields()...,
				)
			}
			if err != nil {
				return nil, 0, err
			}

			req.Latitude = &lat
			req.Longitude = &lon
		}
	}

	// All Pharmacies Fetching
	return u.pr.GetAll(
		ctx,
		&models.GetAllPharmacies{
			Search:    search,
			RadiusM:   req.RadiusM,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,

			ByLocation: req.ByLocation,
			DefaultSorters: utils.DefaultSorters{
				ByCreatedAt: req.ByCreatedAt,
				ByUpdatedAt: req.ByUpdatedAt,
			},

			OffsetPagination: utils.OffsetPagination{
				Page:     page,
				PageSize: pageSize,
			},
		},
	)
}

func (u *pharmacyUsecase) UpdatePharmacy(ctx context.Context, req *models.UpdatePharmacyReq) (*models.Pharmacy, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "pharmacy.usecase.UpdatePharmacy")
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
	name := utils.TrimSpacePtr(req.Name)
	legalName := utils.TrimSpacePtr(req.LegalName)
	description := utils.TrimSpacePtr(req.Description)
	country := utils.ToLowerPtr(utils.TrimSpacePtr(req.Country))
	subdivision1 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision1))
	subdivision2 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision2))
	subdivision3 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision3))
	subdivision4 := utils.ToLowerPtr(utils.TrimSpacePtr(req.Subdivision4))
	street := utils.TrimSpacePtr(req.Street)
	postalCode := utils.ToLowerPtr(utils.TrimSpacePtr(req.PostalCode))
	email := utils.ToLowerPtr(utils.TrimSpacePtr(req.Email))
	phone := utils.TrimSpacePtr(req.Phone)
	website := utils.TrimSpacePtr(req.Website)
	whatsapp := utils.TrimSpacePtr(req.Whatsapp)

	// Data Validation
	if name != nil {
		if ok, why := u.validator.Name(*name, "Name"); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if legalName != nil {
		if ok, why := u.validator.Name(*legalName, "Legal name"); !ok {
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
	if country != nil {
		if ok, why := u.validator.Country(*country, "Country"); !ok {
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
		if ok, why := u.validator.Name(*street, "Street name"); !ok {
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

	// Pharmacy Update
	p, err := u.pr.Update(
		ctx,
		authCtx.AuthID,
		&models.UpdatePharmacy{
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

func (u *pharmacyUsecase) UpdateProfilePicture(ctx context.Context, profilePictureURL string) (*models.Pharmacy, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "pharmacy.usecase.UpdateProfilePicture")
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
	url := strings.TrimSpace(profilePictureURL)

	// Data Validation
	if ok, why := u.validator.StorageURL(url); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}

	// URL Public ID Extraction
	publicID, err := u.storage.ExtractPublicID(url, constants.StorageFilePathTempImages)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeInvalidPayload,
			ce.MsgInvalidPayload,
			err,
			authFields...,
		)
	}

	// Image File Renaming
	res, err := u.storage.Rename(
		ctx,
		publicID,
		"pharmacies/profile_pictures/"+authCtx.PharmacyID.String(),
		true,
		true,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeStorageFileRenamingFailed,
			ce.MsgInternalServer,
			err,
			authFields...,
		)
	}

	// Pharmacy Update
	p, updateErr := u.pr.Update(
		ctx,
		authCtx.AuthID,
		&models.UpdatePharmacy{
			ProfilePicture: &res.SecureURL,
		},
	)
	if updateErr != nil {
		return nil, updateErr.Append(authFields...)
	}

	return p, nil
}

func (u *pharmacyUsecase) UpdateProfileBanner(ctx context.Context, profileBannerURL string) (*models.Pharmacy, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "pharmacy.usecase.UpdateProfileBanner")
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
	url := strings.TrimSpace(profileBannerURL)

	// Data Validation
	if ok, why := u.validator.StorageURL(url); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}

	// URL Public ID Extraction
	publicID, err := u.storage.ExtractPublicID(url, constants.StorageFilePathTempImages)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeInvalidPayload,
			ce.MsgInvalidPayload,
			err,
			authFields...,
		)
	}

	// Image File Renaming
	res, err := u.storage.Rename(
		ctx,
		publicID,
		"pharmacies/profile_banners/"+authCtx.PharmacyID.String(),
		true,
		true,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeStorageFileRenamingFailed,
			ce.MsgInternalServer,
			err,
			authFields...,
		)
	}

	// Pharmacy Update
	p, updateErr := u.pr.Update(
		ctx,
		authCtx.AuthID,
		&models.UpdatePharmacy{
			ProfileBanner: &res.SecureURL,
		},
	)
	if updateErr != nil {
		return nil, updateErr.Append(authFields...)
	}

	return p, nil
}
