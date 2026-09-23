package models

import (
	"encoding/json"

	"github.com/ritchieridanko/apotekly/services/shared/utils"
)

type (
	CreatePharmacyReq struct {
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

	GetAllPharmaciesReq struct {
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

	UpdatePharmacyReq struct {
		Name         *string
		LegalName    *string
		Description  *string
		OnlineHours  *json.RawMessage
		Country      *string
		Subdivision1 *string
		Subdivision2 *string
		Subdivision3 *string
		Subdivision4 *string
		Street       *string
		PostalCode   *string
		Latitude     *float64
		Longitude    *float64
		Email        *string
		Phone        *string
		Website      *string
		Whatsapp     *string
	}
)

func (r *GetAllPharmaciesReq) RequireLocation() bool {
	return r.RadiusM != nil || r.ByLocation != nil
}
