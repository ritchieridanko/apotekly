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
	"github.com/ritchieridanko/apotekly/services/shared/infra/storage"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"go.opentelemetry.io/otel"
)

type ProductImageUsecase interface {
	CreateProductImages(ctx context.Context, req *models.CreateProductImagesReq) (pis []models.ProductImage, err *ce.Error)
	UploadProductImages(ctx context.Context, req *models.UploadProductImagesReq) (pis []models.ProductImage, err *ce.Error)
	GetAllProductImages(ctx context.Context, productID uuid.UUID) (pis []models.ProductImage, err *ce.Error)
	UpdateProductImage(ctx context.Context, req *models.UpdateProductImageReq) (pi *models.ProductImage, err *ce.Error)
	DeleteProductImage(ctx context.Context, req *models.DeleteProductImageReq) (err *ce.Error)
	ReorderProductImage(ctx context.Context, req *models.ReorderProductImageReq) (pi *models.ProductImage, err *ce.Error)
}

type productImageUsecase struct {
	appName        string
	maxTotalImages int64
	pr             repositories.ProductRepository
	pir            repositories.ProductImageRepository
	pc             clients.PharmacyClient
	transactor     *database.Transactor
	storage        *storage.Storage
	validator      *validator.Validator
	logger         *logger.Logger
}

func NewProductImageUsecase(
	appName string,
	maxTotalImages int64,
	pr repositories.ProductRepository,
	pir repositories.ProductImageRepository,
	pc clients.PharmacyClient,
	tx *database.Transactor,
	s *storage.Storage,
	v *validator.Validator,
	l *logger.Logger,
) ProductImageUsecase {
	return &productImageUsecase{
		appName:        appName,
		maxTotalImages: maxTotalImages,
		pr:             pr,
		pir:            pir,
		pc:             pc,
		transactor:     tx,
		storage:        s,
		validator:      v,
		logger:         l,
	}
}

func (u *productImageUsecase) CreateProductImages(ctx context.Context, req *models.CreateProductImagesReq) ([]models.ProductImage, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "productImage.usecase.CreateProductImages")
	defer span.End()

	// Images Total Length Validation
	if int64(len(req.Images)) > u.maxTotalImages {
		return nil, ce.NewError(ce.CodeMaxTotalImagesExceeded, ce.MsgMaxTotalImagesExceeded, nil)
	}

	for i, image := range req.Images {
		// Data Normalization
		name := strings.TrimSpace(image.Name)
		url := strings.TrimSpace(image.URL)
		altText := utils.TrimSpacePtr(image.AltText)

		// Data Validation
		if ok, why := u.validator.Name(name, "Image name"); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
		}
		if ok, why := u.validator.StorageURL(url); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
		}
		if altText != nil {
			if ok, why := u.validator.ImageAltText(*altText); !ok {
				return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
			}
		}

		req.Images[i].Name = name
		req.Images[i].URL = url
		req.Images[i].AltText = altText
	}

	sortOrder := int16(0)
	images := make([]models.CreateProductImage, 0, len(req.Images))
	for _, image := range req.Images {
		// URL Public ID Extraction
		publicID, err := u.storage.ExtractPublicID(image.URL, constants.StorageFilePathTempImages)
		if err != nil {
			return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		}

		// Image File Renaming
		productImageID := utils.MustGenerateUUIDv7()
		res, err := u.storage.Rename(
			ctx,
			publicID,
			"products/"+req.ProductID.String()+"/"+productImageID.String(),
			true,
			true,
		)
		if err != nil {
			return nil, ce.NewError(ce.CodeStorageFileRenamingFailed, ce.MsgInternalServer, err)
		}

		sortOrder++
		images = append(images,
			models.CreateProductImage{
				ID:        productImageID,
				Name:      image.Name,
				URL:       res.SecureURL,
				AltText:   image.AltText,
				SortOrder: sortOrder,
			},
		)
	}

	// Product Images Creation
	return u.pir.CreateMany(
		ctx,
		&models.CreateProductImages{
			ProductID: req.ProductID,
			Images:    images,
		},
	)
}

func (u *productImageUsecase) UploadProductImages(ctx context.Context, req *models.UploadProductImagesReq) ([]models.ProductImage, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "productImage.usecase.UploadProductImages")
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

	// Pharmacy ID Check
	pharmacyID, err := u.pr.GetPharmacyID(ctx, req.ProductID)
	if err != nil {
		return nil, err.Append(authFields...)
	}
	if *authCtx.PharmacyID != pharmacyID {
		return nil, ce.NewError(
			ce.CodeProductPharmacyIDMismatch,
			ce.MsgProductNotFound,
			nil,
			append(
				authFields,
				logger.NewField("product_pharmacy_id", pharmacyID.String()),
			)...,
		)
	}

	// Images Total Length Validation
	totalImages := len(req.Images)
	if totalImages == 0 {
		return nil, ce.NewError(
			ce.CodeInvalidPayload,
			ce.MsgInvalidPayload,
			nil,
			authFields...,
		)
	}

	var pis []models.ProductImage
	err = u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Product Images Total Count
		total, highestSortOrder, err := u.pir.CountAll(ctx, req.ProductID)
		if err != nil {
			return err
		}
		if (total + int64(totalImages)) > u.maxTotalImages {
			return ce.NewError(ce.CodeMaxTotalImagesExceeded, ce.MsgMaxTotalImagesExceeded, nil)
		}

		for i, image := range req.Images {
			// Data Normalization
			name := strings.TrimSpace(image.Name)
			url := strings.TrimSpace(image.URL)
			altText := utils.TrimSpacePtr(image.AltText)

			// Data Validation
			if ok, why := u.validator.Name(name, "Image name"); !ok {
				return ce.NewError(ce.CodeInvalidPayload, why, nil)
			}
			if ok, why := u.validator.StorageURL(url); !ok {
				return ce.NewError(ce.CodeInvalidPayload, why, nil)
			}
			if altText != nil {
				if ok, why := u.validator.ImageAltText(*altText); !ok {
					return ce.NewError(ce.CodeInvalidPayload, why, nil)
				}
			}

			req.Images[i].Name = name
			req.Images[i].URL = url
			req.Images[i].AltText = altText
		}

		images := make([]models.CreateProductImage, 0, totalImages)
		for _, image := range req.Images {
			// URL Public ID Extraction
			publicID, err := u.storage.ExtractPublicID(image.URL, constants.StorageFilePathTempImages)
			if err != nil {
				return ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
			}

			// Image File Renaming
			productImageID := utils.MustGenerateUUIDv7()
			res, err := u.storage.Rename(
				ctx,
				publicID,
				"products/"+req.ProductID.String()+"/"+productImageID.String(),
				true,
				true,
			)
			if err != nil {
				return ce.NewError(ce.CodeStorageFileRenamingFailed, ce.MsgInternalServer, err)
			}

			highestSortOrder++
			images = append(images,
				models.CreateProductImage{
					ID:        productImageID,
					Name:      image.Name,
					URL:       res.SecureURL,
					AltText:   image.AltText,
					SortOrder: highestSortOrder,
				},
			)
		}

		// Product Images Creation
		pis, err = u.pir.CreateMany(
			ctx,
			&models.CreateProductImages{
				ProductID: req.ProductID,
				Images:    images,
			},
		)
		return err
	})
	if err != nil {
		return nil, err.Append(authFields...)
	}

	return pis, nil
}

func (u *productImageUsecase) GetAllProductImages(ctx context.Context, productID uuid.UUID) ([]models.ProductImage, *ce.Error) {
	return u.pir.GetAll(ctx, productID)
}

func (u *productImageUsecase) UpdateProductImage(ctx context.Context, req *models.UpdateProductImageReq) (*models.ProductImage, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "productImage.usecase.UpdateProductImage")
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

	// Pharmacy ID Check
	pharmacyID, err := u.pr.GetPharmacyID(ctx, req.ProductID)
	if err != nil {
		return nil, err.Append(authFields...)
	}
	if *authCtx.PharmacyID != pharmacyID {
		return nil, ce.NewError(
			ce.CodeProductPharmacyIDMismatch,
			ce.MsgProductNotFound,
			nil,
			append(
				authFields,
				logger.NewField("product_pharmacy_id", pharmacyID.String()),
			)...,
		)
	}

	// Data Normalization
	name := utils.TrimSpacePtr(req.Name)
	altText := utils.TrimSpacePtr(req.AltText)

	// Data Validation
	if name != nil {
		if ok, why := u.validator.Name(*name, "Image name"); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if altText != nil {
		if ok, why := u.validator.ImageAltText(*altText); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}

	// Product Image Update
	pi, err := u.pir.Update(
		ctx,
		&models.UpdateProductImageP{
			ProductImageID: req.ProductImageID,
			ProductID:      req.ProductID,
		},
		&models.UpdateProductImageD{
			Name:    name,
			AltText: altText,
		},
	)
	if err != nil {
		return nil, err.Append(authFields...)
	}

	return pi, nil
}

func (u *productImageUsecase) DeleteProductImage(ctx context.Context, req *models.DeleteProductImageReq) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "productImage.usecase.DeleteProductImage")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}
	if authCtx.PharmacyID == nil {
		return ce.NewError(
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
		return ce.NewError(
			ce.CodeRoleMismatch,
			ce.MsgInternalServer,
			err.Unwrap(),
			err.Append(
				authFields...,
			).Fields()...,
		)
	}
	if err != nil {
		return err.Append(authFields...)
	}
	if status != constants.PharmacyStatusActive {
		return ce.NewError(
			ce.CodePharmacyNotActive,
			"Pharmacy account is "+status,
			nil,
			authFields...,
		)
	}

	// Pharmacy ID Check
	pharmacyID, err := u.pr.GetPharmacyID(ctx, req.ProductID)
	if err != nil {
		return err.Append(authFields...)
	}
	if *authCtx.PharmacyID != pharmacyID {
		return ce.NewError(
			ce.CodeProductPharmacyIDMismatch,
			ce.MsgProductNotFound,
			nil,
			append(
				authFields,
				logger.NewField("product_pharmacy_id", pharmacyID.String()),
			)...,
		)
	}

	err = u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Product Image Deletion
		image, err := u.pir.Delete(
			ctx,
			&models.DeleteProductImage{
				ProductImageID: req.ProductImageID,
				ProductID:      req.ProductID,
			},
		)
		if err != nil {
			return err
		}

		// Image File Deletion
		_, delErr := u.storage.Delete(ctx, image, true)
		if delErr != nil {
			return ce.NewError(ce.CodeStorageFileDeletionFailed, ce.MsgInternalServer, delErr)
		}

		return nil
	})
	if err != nil {
		return err.Append(authFields...)
	}

	return nil
}

func (u *productImageUsecase) ReorderProductImage(ctx context.Context, req *models.ReorderProductImageReq) (*models.ProductImage, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "productImage.usecase.ReorderProductImage")
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

	// Pharmacy ID Check
	pharmacyID, err := u.pr.GetPharmacyID(ctx, req.ProductID)
	if err != nil {
		return nil, err.Append(authFields...)
	}
	if *authCtx.PharmacyID != pharmacyID {
		return nil, ce.NewError(
			ce.CodeProductPharmacyIDMismatch,
			ce.MsgProductNotFound,
			nil,
			append(
				authFields,
				logger.NewField("product_pharmacy_id", pharmacyID.String()),
			)...,
		)
	}

	var pi *models.ProductImage
	err = u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Product Images Total Count
		_, highestSortOrder, err := u.pir.CountAll(ctx, req.ProductID)
		if err != nil {
			return err
		}

		// Data Validation
		if ok, why := u.validator.SortOrder(req.NewSortOrder, highestSortOrder); !ok {
			return ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}

		// Product Images Reordering
		pi, err = u.pir.Reorder(
			ctx,
			&models.ReorderProductImage{
				ProductImageID: req.ProductImageID,
				ProductID:      req.ProductID,
			},
			req.NewSortOrder,
		)
		return err
	})
	if err != nil {
		return nil, err.Append(authFields...)
	}

	return pi, nil
}
