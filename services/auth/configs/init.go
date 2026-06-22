package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	cfg "github.com/ritchieridanko/apotekly/services/shared/configs"
	"github.com/spf13/viper"
)

type Config struct {
	App      cfg.App        `mapstructure:"app"`
	Auth     Auth           `mapstructure:"auth"`
	Server   cfg.GRPCServer `mapstructure:"server"`
	Database cfg.Database   `mapstructure:"database"`
	Cache    cfg.Cache      `mapstructure:"cache"`
	Tracer   cfg.Tracer     `mapstructure:"tracer"`
	Broker   Broker         `mapstructure:"broker"`
}

type Auth struct {
	BCrypt cfg.BCrypt `mapstructure:"bcrypt"`
	JWT    cfg.JWT    `mapstructure:"jwt"`

	Duration struct {
		Session      time.Duration `mapstructure:"session"`
		Verification time.Duration `mapstructure:"verification"`
	} `mapstructure:"duration"`
}

type Broker struct {
	Brokers string `mapstructure:"brokers"`

	// Publishers
	AC cfg.Publisher `mapstructure:"ac"`
}

func Init(path string) (*Config, error) {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if env != "dev" && env != "prod" {
		env = "dev"
	}
	if path == "" {
		path = "./configs"
	}

	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	cfg.App.Env = env
	cfg.Server.Addr = cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	cfg.Cache.Addr = cfg.Cache.Host + ":" + strconv.Itoa(cfg.Cache.Port)
	cfg.Tracer.Addr = cfg.Tracer.Host + ":" + strconv.Itoa(cfg.Tracer.Port)
	cfg.Database.DSN = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	return &cfg, nil
}
