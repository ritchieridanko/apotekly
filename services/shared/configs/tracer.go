package configs

type Tracer struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
