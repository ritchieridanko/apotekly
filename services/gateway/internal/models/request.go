package models

import (
	"encoding/json"
	"time"
)

// Auth Service
type (
	ChangeEmailReq struct {
		Password string
		NewEmail string
	}

	ChangePasswordReq struct {
		OldPassword string
		NewPassword string
	}

	ConfirmPasswordResetReq struct {
		PasswordResetToken string
		NewPassword        string
	}

	SignInReq struct {
		Email    string
		Password string
	}

	SignUpReq struct {
		Email    string
		Password string
	}

	VerifyEmailReq struct {
		RefreshToken      string
		VerificationToken string
	}
)

// Pharmacy Service
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

// Product Service
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

// User Service
type (
	CreateUserReq struct {
		Name      string
		Sex       *string
		Birthdate *time.Time
		Phone     *string
	}

	UpdateUserReq struct {
		Name      *string
		Sex       *string
		Birthdate *time.Time
		Phone     *string
	}

	// User Address Service
	CreateAddressReq struct {
		Label        string
		Recipient    string
		Phone        string
		Notes        *string
		IsPrimary    bool
		Country      string
		Subdivision1 *string
		Subdivision2 *string
		Subdivision3 *string
		Subdivision4 *string
		Street       string
		PostalCode   string
		Latitude     float64
		Longitude    float64
	}

	GetAllAddressesReq struct {
		Page     int32
		PageSize int32
	}

	UpdateAddressReq struct {
		// Params
		AddressID uint64

		// Data
		Label        *string
		Recipient    *string
		Phone        *string
		Notes        *string
		Country      *string
		Subdivision1 *string
		Subdivision2 *string
		Subdivision3 *string
		Subdivision4 *string
		Street       *string
		PostalCode   *string
		Latitude     *float64
		Longitude    *float64
	}
)
