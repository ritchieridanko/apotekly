package configs

type Service struct {
	Name string `mapstructure:"name"`
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
