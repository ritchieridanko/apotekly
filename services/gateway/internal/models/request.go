package models

import "time"

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
)
