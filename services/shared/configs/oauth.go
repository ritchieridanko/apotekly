package configs

type OAuth struct {
	Name        string `mapstructure:"name"`
	ClientID    string `mapstructure:"client_id"`
	Secret      string `mapstructure:"secret"`
	RedirectURL string `mapstructure:"redirect_url"`
	OIDCIssuer  string `mapstructure:"oidc_issuer"`
	Scopes      string `mapstructure:"scopes"`
}
