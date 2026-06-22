package models

type (
	CreateSessionReq struct {
		AuthID          uint64
		Role            string
		IsEmailVerified bool
	}

	SignUpReq struct {
		Email    string
		Password string
	}
)
