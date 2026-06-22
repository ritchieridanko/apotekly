package models

import "time"

type (
	AccessToken struct {
		Token            string
		ExpiresInSeconds uint64
	}

	RefreshToken struct {
		Token            string
		ExpiresInSeconds uint64
	}

	AuthToken struct {
		AccessToken  *AccessToken
		RefreshToken *RefreshToken
	}

	CreateVerificationToken struct {
		AuthID   uint64
		Token    string
		Duration time.Duration
	}
)
