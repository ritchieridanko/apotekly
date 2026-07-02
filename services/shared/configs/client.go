package configs

type Client struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
