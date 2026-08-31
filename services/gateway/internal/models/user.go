package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Name           string
	Sex            *string
	Birthdate      *time.Time
	Phone          *string
	ProfilePicture *string
	ProfileBanner  *string
}
