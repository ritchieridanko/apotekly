package configs

type App struct {
	Name    string `mapstructure:"name"`
	LogoURL string `mapstructure:"logo_url"`
	Env     string
}
