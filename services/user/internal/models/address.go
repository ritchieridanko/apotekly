package models

import "time"

type (
	Address struct {
		ID           uint64
		AuthID       uint64
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
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	CreateAddress struct {
		AuthID       uint64
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
