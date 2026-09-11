package policies

import "github.com/ritchieridanko/apotekly/services/shared/constants"

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
	"/pharmacy.v1.PharmacyService/CreatePharmacy": {
		authenticated: true,
		verified:      true,
		roles: map[string]struct{}{
			constants.RoleUser: {},
		},
	},
}
