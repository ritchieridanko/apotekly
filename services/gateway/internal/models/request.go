package models

type (
	SignInReq struct {
		Email    string
		Password string
	}

	SignUpReq struct {
		Email    string
		Password string
	}

	VerifyEmailReq struct {
		RefreshToken      string
		VerificationToken string
	}
)
