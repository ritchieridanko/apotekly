package models

import "github.com/google/uuid"

type (
	EmailChangeEmail struct {
		EventID  uuid.UUID
		OldEmail string
		NewEmail string
		Role     string
		Token    string
	}

	VerificationEmail struct {
		Recipient string
		Role      string
		Token     string
	}

	WelcomeEmail struct {
		Recipient         string
		Role              string
		IsEmailVerified   bool
		VerificationToken string
	}
)
