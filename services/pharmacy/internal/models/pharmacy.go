package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
)

type (
	Pharmacy struct {
		ID             uuid.UUID
		Name           string
		LegalName      *string
		Description    *string
		Status         string
		OnlineHours    *json.RawMessage
		Country        string
		Subdivision1   *string
		Subdivision2   *string
		Subdivision3   *string
		Subdivision4   *string
		Street         string
		PostalCode     string
		Latitude       float64
		Longitude      float64
		Email          *string
		Phone          *string
		Website        *string
		Whatsapp       *string
		ProfilePicture *string
		ProfileBanner  *string
		VerifiedAt     *time.Time
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	PharmacySummary struct {
		ID             uuid.UUID
		Name           string
		LegalName      *string
		OnlineHours    *json.RawMessage
		ProfilePicture *string
		DistanceM      *float64
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	CreatePharmacy struct {
		ID           uuid.UUID
		AuthID       uint64
		Name         string
		LegalName    *string
		Description  *string
		OnlineHours  *json.RawMessage
		Country      string
		Subdivision1 *string
		Subdivision2 *string
		Subdivision3 *string
		Subdivision4 *string
		Street       string
		PostalCode   string
		Latitude     float64
		Longitude    float64
		Email        *string
		Phone        *string
		Website      *string
		Whatsapp     *string
	}

	GetAllPharmacies struct {
		// Queries
		Search    *string
		RadiusM   *uint32
		Latitude  *float64
		Longitude *float64

		// Sorters
		ByLocation *utils.Sorter
		utils.DefaultSorters

		// Pagination
		utils.OffsetPagination
	}

	UpdatePharmacy struct {
		Name           *string
		LegalName      *string
		Description    *string
		OnlineHours    *json.RawMessage
		Country        *string
		Subdivision1   *string
		Subdivision2   *string
		Subdivision3   *string
		Subdivision4   *string
		Street         *string
		PostalCode     *string
		Latitude       *float64
		Longitude      *float64
		Email          *string
		Phone          *string
		Website        *string
		Whatsapp       *string
		ProfilePicture *string
		ProfileBanner  *string
	}
)

func (p *GetAllPharmacies) RequireLocation() bool {
	return p.RadiusM != nil || p.ByLocation != nil
}
