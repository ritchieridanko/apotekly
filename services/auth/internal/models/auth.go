package models

import "time"

type (
	Auth struct {
		ID                uint64
		Email             string
		Password          *string
		Role              string
		EmailVerifiedAt   *time.Time
		EmailChangedAt    *time.Time
		PasswordChangedAt *time.Time
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	CreateAuth struct {
		Email           string
		Password        *string
		Role            string
		EmailVerifiedAt *time.Time
	}
)

func (a *Auth) IsEmailVerified() bool {
	return a.EmailVerifiedAt != nil
}
