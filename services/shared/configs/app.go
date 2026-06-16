package configs

type App struct {
	Name string `mapstructure:"name"`
	Env  string
}
