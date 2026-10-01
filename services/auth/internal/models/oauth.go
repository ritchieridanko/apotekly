package models

import "time"

type (
	OAuth struct {
		ID       uint64
		Provider string
		Subject  string
	}

	CreateOAuth struct {
		AuthID   uint64
		Provider string
		Subject  string
	}

	StoreOAuthState struct {
		State    string
		Provider string
		Duration time.Duration
	}
)
