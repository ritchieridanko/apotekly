package configs

import "time"

type Publisher struct {
	Name         string        `mapstructure:"name"`
	Balancer     string        `mapstructure:"balancer"`
	BatchSize    int           `mapstructure:"batch_size"`
	BatchTimeout time.Duration `mapstructure:"batch_timeout"`
}

type Subscriber struct {
	Name           string        `mapstructure:"name"`
	MaxBytes       int           `mapstructure:"max_bytes"`
	MaxWait        time.Duration `mapstructure:"max_wait"`
	CommitInterval time.Duration `mapstructure:"commit_interval"`
}
