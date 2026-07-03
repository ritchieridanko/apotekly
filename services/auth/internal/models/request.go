package models

type (
	CreateSessionReq struct {
		AuthID          uint64
		Role            string
		IsEmailVerified bool
	}

	SignInReq struct {
		Email    string
		Password string
	}

	SignUpReq struct {
		Email    string
		Password string
	}
)
