package dtos

// Requests
type (
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

	IsEmailAvailableResponse struct {
		IsAvailable bool `json:"is_available"`
	}

	SignInResponse struct {
		Auth        *Auth        `json:"auth,omitempty"`
		AccessToken *AccessToken `json:"access_token,omitempty"`
	}

	SignUpResponse struct {
		Auth        *Auth        `json:"auth,omitempty"`
		AccessToken *AccessToken `json:"access_token,omitempty"`
	}
)
