package models

import "github.com/google/uuid"

type (
	ProductImage struct {
		ID        uuid.UUID
		ProductID uuid.UUID
		Name      string
		Image     string
		AltText   *string
		SortOrder int16
	}

	CreateProductImage struct {
		ID        uuid.UUID
		Name      string
		URL       string
		AltText   *string
		SortOrder int16
	}

	CreateProductImages struct {
		ProductID uuid.UUID
		Images    []CreateProductImage
	}

	UpdateProductImageP struct {
		ProductImageID uuid.UUID
		ProductID      uuid.UUID
	}

	UpdateProductImageD struct {
		Name    *string
		AltText *string
	}
)
