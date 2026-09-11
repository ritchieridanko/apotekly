package models

import "encoding/json"

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
)
