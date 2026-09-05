package models

import "time"

// User
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
)

// Address
type (
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
)
