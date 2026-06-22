package configs

import "time"

type Publisher struct {
	Name         string        `mapstructure:"name"`
	Balancer     string        `mapstructure:"balancer"`
	BatchSize    int           `mapstructure:"batch_size"`
	BatchTimeout time.Duration `mapstructure:"batch_timeout"`
}
