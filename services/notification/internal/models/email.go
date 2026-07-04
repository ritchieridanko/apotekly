package models

type (
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
