package models

import "time"

type (
	Session struct {
		ID           uint64
		ParentID     *uint64
		AuthID       uint64
		RefreshToken string
		IPAddress    string
		UserAgent    string
		CreatedAt    time.Time
		ExpiresAt    time.Time
	}

	CreateSession struct {
		ParentID     *uint64
		AuthID       uint64
		RefreshToken string
		IPAddress    string
		UserAgent    string
		ExpiresAt    time.Time
	}

	RevokeActiveSession struct {
		AuthID    uint64
		IPAddress string
		UserAgent string
		ExpiresAt time.Time
	}
)
