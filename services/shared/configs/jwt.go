package configs

import "time"

type JWT struct {
	Issuer   string        `mapstructure:"issuer"`
	Secret   string        `mapstructure:"secret"`
	Duration time.Duration `mapstructure:"duration"`
}
