package models

import "time"

type (
	CreateUserReq struct {
		Name      string
		Sex       *string
		Birthdate *time.Time
		Phone     *string
	}
)
