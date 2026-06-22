package configs

import "time"

type Cache struct {
	Addr            string
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Pass            string        `mapstructure:"pass"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxActiveConns  int           `mapstructure:"max_active_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`

	Timeout struct {
		Dial  time.Duration `mapstructure:"dial"`
		Read  time.Duration `mapstructure:"read"`
		Write time.Duration `mapstructure:"write"`
	} `mapstructure:"timeout"`
}
