package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	cfg "github.com/ritchieridanko/apotekly/services/shared/configs"
	"github.com/spf13/viper"
)

type Config struct {
	App      cfg.App        `mapstructure:"app"`
	Server   cfg.GRPCServer `mapstructure:"server"`
	Service  Service        `mapstructure:"service"`
	Database cfg.Database   `mapstructure:"database"`
	Storage  cfg.Storage    `mapstructure:"storage"`
	Tracer   cfg.Tracer     `mapstructure:"tracer"`
}

type Service struct {
	Auth cfg.Service `mapstructure:"auth"`
	User cfg.Service `mapstructure:"user"`
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
	cfg.Service.Auth.Addr = cfg.Service.Auth.Host + ":" + strconv.Itoa(cfg.Service.Auth.Port)
	cfg.Service.User.Addr = cfg.Service.User.Host + ":" + strconv.Itoa(cfg.Service.User.Port)
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
