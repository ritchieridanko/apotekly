package dtos

// Requests
type (
	ChangeEmailRequest struct {
		Password string `json:"password" binding:"required"`
		NewEmail string `json:"new_email" binding:"required"`
	}

	ChangePasswordRequest struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	IsEmailAvailableRequest struct {
		Email string `form:"email" binding:"required"`
	}

	SignInRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	SignUpRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	VerifyEmailRequest struct {
		VerificationToken string `form:"token" binding:"required"`
	}
)

// Responses
type (
	AccessToken struct {
		Token            string `json:"token"`
		ExpiresInSeconds uint64 `json:"expires_in_seconds"`
	}

	Auth struct {
		Email           string `json:"email"`
		Role            string `json:"role"`
		IsEmailVerified bool   `json:"is_email_verified"`
	}

	ChangeEmailResponse struct {
		Email string `json:"email"`
	}

	IsEmailAvailableResponse struct {
		IsAvailable bool `json:"is_available"`
	}

	ResendVerificationResponse struct {
		Email string `json:"email"`
	}

	RotateAuthTokenResponse struct {
		AccessToken *AccessToken `json:"access_token,omitempty"`
	}

	SignInResponse struct {
		Auth        *Auth        `json:"auth,omitempty"`
		AccessToken *AccessToken `json:"access_token,omitempty"`
	}

	SignUpResponse struct {
		Auth        *Auth        `json:"auth,omitempty"`
		AccessToken *AccessToken `json:"access_token,omitempty"`
	}

	VerifyEmailResponse struct {
		Auth        *Auth        `json:"auth,omitempty"`
		AccessToken *AccessToken `json:"access_token,omitempty"`
	}
)
