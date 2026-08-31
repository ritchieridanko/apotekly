package dtos

import "time"

// Requests
type (
	CreateUserRequest struct {
		Name      string     `json:"name" binding:"required"`
		Sex       *string    `json:"sex"`
		Birthdate *time.Time `json:"birthdate" time_format:"2006-01-02"`
		Phone     *string    `json:"phone"`
	}
)

// Responses
type (
	User struct {
		ID             string     `json:"id"`
		Name           string     `json:"name"`
		Sex            *string    `json:"sex"`
		Birthdate      *time.Time `json:"birthdate"`
		Phone          *string    `json:"phone"`
		ProfilePicture *string    `json:"profile_picture"`
		ProfileBanner  *string    `json:"profile_banner"`
	}

	CreateUserResponse struct {
		User *User `json:"user,omitempty"`
	}

	GetMeResponse struct {
		User *User `json:"user,omitempty"`
	}
)
