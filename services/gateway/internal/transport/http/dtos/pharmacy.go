package dtos

import (
	"encoding/json"
	"time"
)

// Requests
type (
	CreatePharmacyRequest struct {
		// Data
		Name         string           `json:"name" binding:"required"`
		LegalName    *string          `json:"legal_name"`
		Description  *string          `json:"description"`
		OnlineHours  *json.RawMessage `json:"online_hours"`
		Country      string           `json:"country" binding:"required"`
		Subdivision1 *string          `json:"subdivision_1"`
		Subdivision2 *string          `json:"subdivision_2"`
		Subdivision3 *string          `json:"subdivision_3"`
		Subdivision4 *string          `json:"subdivision_4"`
		Street       string           `json:"street" binding:"required"`
		PostalCode   string           `json:"postal_code" binding:"required"`
		Latitude     float64          `json:"latitude" binding:"required"`
		Longitude    float64          `json:"longitude" binding:"required"`
		Email        *string          `json:"email"`
		Phone        *string          `json:"phone"`
		Website      *string          `json:"website"`
		Whatsapp     *string          `json:"whatsapp"`

		// Auth Info
		RememberMe bool `json:"remember_me"`
	}

	UpdatePharmacyRequest struct {
		Name         *string          `json:"name"`
		LegalName    *string          `json:"legal_name"`
		Description  *string          `json:"description"`
		OnlineHours  *json.RawMessage `json:"online_hours"`
		Country      *string          `json:"country"`
		Subdivision1 *string          `json:"subdivision_1"`
		Subdivision2 *string          `json:"subdivision_2"`
		Subdivision3 *string          `json:"subdivision_3"`
		Subdivision4 *string          `json:"subdivision_4"`
		Street       *string          `json:"street"`
		PostalCode   *string          `json:"postal_code"`
		Latitude     *float64         `json:"latitude"`
		Longitude    *float64         `json:"longitude"`
		Email        *string          `json:"email"`
		Phone        *string          `json:"phone"`
		Website      *string          `json:"website"`
		Whatsapp     *string          `json:"whatsapp"`
	}

	PharmacyUpdateProfilePictureRequest struct {
		ProfilePictureURL string `json:"profile_picture_url" binding:"required"`
	}

	PharmacyUpdateProfileBannerRequest struct {
		ProfileBannerURL string `json:"profile_banner_url" binding:"required"`
	}
)

// Responses
type (
	Pharmacy struct {
		ID             string           `json:"id"`
		Name           string           `json:"name"`
		LegalName      *string          `json:"legal_name"`
		Description    *string          `json:"description"`
		Status         string           `json:"status"`
		OnlineHours    *json.RawMessage `json:"online_hours"`
		Country        string           `json:"country"`
		Subdivision1   *string          `json:"subdivision_1"`
		Subdivision2   *string          `json:"subdivision_2"`
		Subdivision3   *string          `json:"subdivision_3"`
		Subdivision4   *string          `json:"subdivision_4"`
		Street         string           `json:"street"`
		PostalCode     string           `json:"postal_code"`
		Latitude       float64          `json:"latitude"`
		Longitude      float64          `json:"longitude"`
		Email          *string          `json:"email"`
		Phone          *string          `json:"phone"`
		Website        *string          `json:"website"`
		Whatsapp       *string          `json:"whatsapp"`
		ProfilePicture *string          `json:"profile_picture"`
		ProfileBanner  *string          `json:"profile_banner"`
		VerifiedAt     *time.Time       `json:"verified_at"`
		CreatedAt      *time.Time       `json:"created_at"`
		UpdatedAt      *time.Time       `json:"updated_at"`
	}

	CreatePharmacyResponse struct {
		Pharmacy    *Pharmacy    `json:"pharmacy,omitempty"`
		AccessToken *AccessToken `json:"access_token,omitempty"`
	}

	PharmacyGetMeResponse struct {
		Pharmacy *Pharmacy `json:"pharmacy,omitempty"`
	}

	UpdatePharmacyResponse struct {
		Pharmacy *Pharmacy `json:"pharmacy,omitempty"`
	}

	PharmacyUpdateProfilePictureResponse struct {
		Pharmacy *Pharmacy `json:"pharmacy,omitempty"`
	}

	PharmacyUpdateProfileBannerResponse struct {
		Pharmacy *Pharmacy `json:"pharmacy,omitempty"`
	}
)
