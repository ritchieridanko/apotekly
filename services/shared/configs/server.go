package configs

import "time"

type GRPCServer struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	Timeout struct {
		Shutdown time.Duration `mapstructure:"shutdown"`
	} `mapstructure:"timeout"`
}
