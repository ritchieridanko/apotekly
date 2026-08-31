package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID             uuid.UUID
		AuthID         uint64
		Name           string
		Sex            *string
		Birthdate      *time.Time
		Phone          *string
		ProfilePicture *string
		ProfileBanner  *string
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	CreateUser struct {
		ID        uuid.UUID
		AuthID    uint64
		Name      string
		Sex       *string
		Birthdate *time.Time
		Phone     *string
	}

	UpdateUser struct {
		Name      *string
		Sex       *string
		Birthdate *time.Time
		Phone     *string
	}
)
