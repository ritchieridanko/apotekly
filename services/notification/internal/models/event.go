package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	Event struct {
		ID      uuid.UUID
		Topic   string
		Payload json.RawMessage
	}

	EventInbox struct {
		ID              uuid.UUID
		Topic           string
		Payload         json.RawMessage
		ReceivedAt      time.Time
		CompletedAt     *time.Time
		RetryCount      int
		NextRetryAt     time.Time
		LastAttemptedAt *time.Time
		LastError       *string
	}

	CreateEventInbox struct {
		ID      uuid.UUID
		Topic   string
		Payload json.RawMessage
	}

	RetryEventInbox struct {
		NextRetryAt time.Time
		Error       string
	}
)

type (
	EventAC struct {
		ID                uuid.UUID  `json:"id"`
		AuthID            uint64     `json:"auth_id"`
		Email             string     `json:"email"`
		Role              string     `json:"role"`
		EmailVerifiedAt   *time.Time `json:"email_verified_at"`
		Session           *string    `json:"session"`
		VerificationToken *string    `json:"verification_token"`
		CreatedAt         *time.Time `json:"created_at"`
	}

	EventAEVR struct {
		ID        uuid.UUID  `json:"id"`
		AuthID    uint64     `json:"auth_id"`
		Email     string     `json:"email"`
		Role      string     `json:"role"`
		Token     string     `json:"token"`
		CreatedAt *time.Time `json:"created_at"`
	}
)
