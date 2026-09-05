package policies

type AuthPolicy struct {
	authenticated bool
	verified      bool
	roles         map[string]struct{}
}

func (p *AuthPolicy) MustBeAuthenticated() bool {
	return p.authenticated
}

func (p *AuthPolicy) MustBeVerified() bool {
	return p.verified
}

func (p *AuthPolicy) RequireRole() bool {
	return len(p.roles) > 0
}

func (p *AuthPolicy) IsRoleAuthorized(role string) bool {
	_, exists := p.roles[role]
	return exists
}

var AuthPolicies map[string]AuthPolicy = map[string]AuthPolicy{
	// User
	"/user.v1.UserService/CreateUser": {authenticated: true},
	"/user.v1.UserService/GetMe":      {authenticated: true},
	"/user.v1.UserService/UpdateUser": {authenticated: true},

	// Address
	"/user.v1.AddressService/CreateAddress": {authenticated: true},
}
