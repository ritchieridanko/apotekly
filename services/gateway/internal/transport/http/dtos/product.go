package dtos

import "time"

// Requests
type (
	CreateProductRequest struct {
		BrandName      string                `json:"brand_name" binding:"required"`
		GenericName    string                `json:"generic_name" binding:"required"`
		Description    *string               `json:"description"`
		RequiresRX     bool                  `json:"requires_rx"`
		DosageForm     string                `json:"dosage_form" binding:"required"`
		Strength       *string               `json:"strength"`
		PackUnit       string                `json:"pack_unit" binding:"required"`
		PackSize       int                   `json:"pack_size" binding:"required"`
		HeightCM       float64               `json:"height_cm" binding:"required"`
		LengthCM       float64               `json:"length_cm" binding:"required"`
		WidthCM        float64               `json:"width_cm" binding:"required"`
		WeightG        float64               `json:"weight_g" binding:"required"`
		Price          string                `json:"price" binding:"required"`
		Currency       string                `json:"currency" binding:"required"`
		Quantity       int                   `json:"quantity" binding:"required"`
		IsActive       bool                  `json:"is_active"`
		ManufacturedBy string                `json:"manufactured_by" binding:"required"`
		ManufacturedIn string                `json:"manufactured_in" binding:"required"`
		RegAuthority   string                `json:"reg_authority" binding:"required"`
		RegIdentifier  string                `json:"reg_identifier" binding:"required"`
		Images         []ProductImageRequest `json:"images"`
	}

	// Product Image
	ProductImageRequest struct {
		Name    string  `json:"name" binding:"required"`
		URL     string  `json:"url" binding:"required"`
		AltText *string `json:"alt_text"`
	}

	UploadProductImagesRequest struct {
		Images []ProductImageRequest `json:"images"`
	}

	UpdateProductImageRequest struct {
		Name    *string `json:"name"`
		AltText *string `json:"alt_text"`
	}

	ReorderProductImageRequest struct {
		NewSortOrder int16 `json:"new_sort_order" binding:"required"`
	}
)

// Responses
type (
	Product struct {
		ID             string         `json:"id"`
		PharmacyID     string         `json:"pharmacy_id"`
		BrandName      string         `json:"brand_name"`
		GenericName    string         `json:"generic_name"`
		Description    *string        `json:"description"`
		RequiresRX     bool           `json:"requires_rx"`
		DosageForm     string         `json:"dosage_form"`
		Strength       *string        `json:"strength"`
		PackUnit       string         `json:"pack_unit"`
		PackSize       int            `json:"pack_size"`
		HeightCM       float64        `json:"height_cm"`
		LengthCM       float64        `json:"length_cm"`
		WidthCM        float64        `json:"width_cm"`
		WeightG        float64        `json:"weight_g"`
		Price          string         `json:"price"`
		Currency       string         `json:"currency"`
		Quantity       int            `json:"quantity"`
		IsActive       bool           `json:"is_active"`
		ManufacturedBy string         `json:"manufactured_by"`
		ManufacturedIn string         `json:"manufactured_in"`
		RegAuthority   string         `json:"reg_authority"`
		RegIdentifier  string         `json:"reg_identifier"`
		Images         []ProductImage `json:"images"`
		CreatedAt      *time.Time     `json:"created_at"`
		UpdatedAt      *time.Time     `json:"updated_at"`
	}

	CreateProductResponse struct {
		Product *Product `json:"product,omitempty"`
	}

	GetProductByIDResponse struct {
		Product *Product `json:"product,omitempty"`
	}

	// Product Image
	ProductImage struct {
		ID        string  `json:"id"`
		ProductID string  `json:"product_id"`
		Name      string  `json:"name"`
		Image     string  `json:"image"`
		AltText   *string `json:"alt_text"`
		SortOrder int16   `json:"sort_order"`
	}

	UploadProductImagesResponse struct {
		ProductImages []ProductImage `json:"product_images"`
	}

	GetAllProductImagesResponse struct {
		ProductImages []ProductImage `json:"product_images"`
	}

	UpdateProductImageResponse struct {
		ProductImage *ProductImage `json:"product_image,omitempty"`
	}

	ReorderProductImageResponse struct {
		ProductImage *ProductImage `json:"product_image,omitempty"`
	}
)
