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
)
