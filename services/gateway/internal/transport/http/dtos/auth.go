package dtos

// Requests
type (
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

	SignUpResponse struct {
		Auth        *Auth        `json:"auth,omitempty"`
		AccessToken *AccessToken `json:"access_token,omitempty"`
	}
)
