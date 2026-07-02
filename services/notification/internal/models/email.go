package models

type (
	WelcomeEmail struct {
		Recipient         string
		Role              string
		IsEmailVerified   bool
		VerificationToken string
	}
)
