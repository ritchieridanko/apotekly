package policies

import "github.com/ritchieridanko/apotekly/services/shared/constants"

type AuthPolicy struct {
	authenticated bool
	authOptional  bool
	verified      bool
	roles         map[string]struct{}
}

func (p *AuthPolicy) MustBeAuthenticated() bool {
	return p.authenticated
}

func (p *AuthPolicy) IsAuthOptional() bool {
	return p.authOptional
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
	// Product
	"/product.v1.ProductService/CreateProduct": {
		authenticated: true,
		roles: map[string]struct{}{
			constants.RolePharmacy: {},
		},
	},

	// Product Image
	"/product.v1.ProductImageService/UploadProductImages": {
		authenticated: true,
		roles: map[string]struct{}{
			constants.RolePharmacy: {},
		},
	},
	"/product.v1.ProductImageService/UpdateProductImage": {
		authenticated: true,
		roles: map[string]struct{}{
			constants.RolePharmacy: {},
		},
	},
	"/product.v1.ProductImageService/DeleteProductImage": {
		authenticated: true,
		roles: map[string]struct{}{
			constants.RolePharmacy: {},
		},
	},
	"/product.v1.ProductImageService/ReorderProductImage": {
		authenticated: true,
		roles: map[string]struct{}{
			constants.RolePharmacy: {},
		},
	},
}
