package models

type (
	Auth struct {
		Email           string
		Role            string
		IsEmailVerified bool
	}
)
