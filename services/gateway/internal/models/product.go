package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Product struct {
		ID             uuid.UUID
		PharmacyID     uuid.UUID
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
		CreatedAt      *time.Time
		UpdatedAt      *time.Time
	}
)
