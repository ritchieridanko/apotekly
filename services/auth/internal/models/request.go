package models

type (
	ChangeEmailReq struct {
		Password string
		NewEmail string
	}

	ChangePasswordReq struct {
		OldPassword string
		NewPassword string
	}

	CreateSessionReq struct {
		AuthID          uint64
		Role            string
		IsEmailVerified bool
	}

	RefreshSessionReq struct {
		AuthID          uint64
		Role            string
		IsEmailVerified bool
		RefreshToken    string
	}

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
