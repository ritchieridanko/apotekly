package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type JWT struct {
	issuer   string
	secret   string
	duration time.Duration
}

func Init(issuer, secret string, dn time.Duration) *JWT {
	return &JWT{
		issuer:   issuer,
		secret:   secret,
		duration: dn,
	}
}

func (j *JWT) Generate(authID uint64, role string, isEmailVerified bool, now *time.Time) (string, error) {
	if now == nil {
		t := time.Now().UTC()
		now = &t
	}
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claim{
			AuthID:          authID,
			Role:            role,
			IsEmailVerified: isEmailVerified,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				Subject:   "auth_" + strconv.FormatUint(authID, 10),
				IssuedAt:  &jwt.NumericDate{Time: *now},
				ExpiresAt: &jwt.NumericDate{Time: now.Add(j.duration)},
			},
		},
	).SignedString(
		[]byte(j.secret),
	)
}

func (j *JWT) Parse(token string) (*Claim, error) {
	t, err := jwt.ParseWithClaims(
		token,
		&Claim{},
		func(t *jwt.Token) (any, error) {
			return []byte(j.secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := t.Claims.(*Claim)
	if !ok {
		return nil, ce.ErrInvalidJWTClaim
	}
	return claim, nil
}
