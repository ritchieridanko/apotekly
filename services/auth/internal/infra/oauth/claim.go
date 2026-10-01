package oauth

type Claim struct {
	Subject         string `json:"sub"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"email_verified"`
}
