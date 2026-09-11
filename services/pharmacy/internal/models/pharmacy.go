package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
)
