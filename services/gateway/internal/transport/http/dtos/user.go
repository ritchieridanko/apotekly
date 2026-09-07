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

	UpdateUserRequest struct {
		Name      *string    `json:"name"`
		Sex       *string    `json:"sex"`
		Birthdate *time.Time `json:"birthdate" time_format:"2006-01-02"`
		Phone     *string    `json:"phone"`
	}

	// Address
	CreateAddressRequest struct {
		Label        string  `json:"label" binding:"required"`
		Recipient    string  `json:"recipient" binding:"required"`
		Phone        string  `json:"phone" binding:"required"`
		Notes        *string `json:"notes"`
		IsPrimary    bool    `json:"is_primary"`
		Country      string  `json:"country" binding:"required"`
		Subdivision1 *string `json:"subdivision_1"`
		Subdivision2 *string `json:"subdivision_2"`
		Subdivision3 *string `json:"subdivision_3"`
		Subdivision4 *string `json:"subdivision_4"`
		Street       string  `json:"street" binding:"required"`
		PostalCode   string  `json:"postal_code" binding:"required"`
		Latitude     float64 `json:"latitude" binding:"required"`
		Longitude    float64 `json:"longitude" binding:"required"`
	}

	GetAllAddressesRequest struct {
		PaginationParams
	}

	UpdateAddressRequest struct {
		Label        *string  `json:"label"`
		Recipient    *string  `json:"recipient"`
		Phone        *string  `json:"phone"`
		Notes        *string  `json:"notes"`
		Country      *string  `json:"country"`
		Subdivision1 *string  `json:"subdivision_1"`
		Subdivision2 *string  `json:"subdivision_2"`
		Subdivision3 *string  `json:"subdivision_3"`
		Subdivision4 *string  `json:"subdivision_4"`
		Street       *string  `json:"street"`
		PostalCode   *string  `json:"postal_code"`
		Latitude     *float64 `json:"latitude"`
		Longitude    *float64 `json:"longitude"`
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

	UpdateUserResponse struct {
		User *User `json:"user,omitempty"`
	}

	// Address
	Address struct {
		ID           uint64     `json:"id"`
		Label        string     `json:"label"`
		Recipient    string     `json:"recipient"`
		Phone        string     `json:"phone"`
		Notes        *string    `json:"notes"`
		IsPrimary    bool       `json:"is_primary"`
		Country      string     `json:"country"`
		Subdivision1 *string    `json:"subdivision_1"`
		Subdivision2 *string    `json:"subdivision_2"`
		Subdivision3 *string    `json:"subdivision_3"`
		Subdivision4 *string    `json:"subdivision_4"`
		Street       string     `json:"street"`
		PostalCode   string     `json:"postal_code"`
		Latitude     float64    `json:"latitude"`
		Longitude    float64    `json:"longitude"`
		CreatedAt    *time.Time `json:"created_at"`
		UpdatedAt    *time.Time `json:"updated_at"`
	}

	CreateAddressResponse struct {
		Address *Address `json:"address,omitempty"`
	}

	GetAllAddressesResponse struct {
		Addresses []Address `json:"addresses"`
	}

	UpdateAddressResponse struct {
		Address *Address `json:"address,omitempty"`
	}

	SetPrimaryAddressResponse struct {
		Address *Address `json:"address,omitempty"`
	}
)
