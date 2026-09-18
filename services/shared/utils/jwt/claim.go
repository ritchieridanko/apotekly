package jwt

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	AuthID          uint64
	Role            string
	IsEmailVerified bool
	PharmacyID      *string
	jwt.RegisteredClaims
}
