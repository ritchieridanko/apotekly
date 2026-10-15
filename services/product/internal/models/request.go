package models

import "github.com/google/uuid"

// Product
type (
	CreateProductReq struct {
		BrandName      string
		GenericName    string
		Description    *string
		RequiresRX     bool
		DosageForm     string
		Strength       *string
		PackUnit       string
		PackSize       int
		HeightCM       float64
		LengthCM       float64
		WidthCM        float64
		WeightG        float64
		Price          string
		Currency       string
		Quantity       int
		IsActive       bool
		ManufacturedBy string
		ManufacturedIn string
		RegAuthority   string
		RegIdentifier  string
	}
)

// Product Image
type (
	UploadProductImageReq struct {
		Name    string
		URL     string
		AltText *string
	}

	UploadProductImagesReq struct {
		// Params
		ProductID uuid.UUID

		// Data
		Images []UploadProductImageReq
	}

	UpdateProductImageReq struct {
		// Params
		ProductImageID uuid.UUID
		ProductID      uuid.UUID

		// Data
		Name    *string
		AltText *string
	}

	ReorderProductImageReq struct {
		// Params
		ProductImageID uuid.UUID
		ProductID      uuid.UUID

		// Data
		NewSortOrder int16
	}
)
