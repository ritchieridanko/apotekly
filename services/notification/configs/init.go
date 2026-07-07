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
	App      cfg.App      `mapstructure:"app"`
	Client   cfg.Client   `mapstructure:"client"`
	Database cfg.Database `mapstructure:"database"`
	Mailer   cfg.Mailer   `mapstructure:"mailer"`
	Tracer   cfg.Tracer   `mapstructure:"tracer"`
	Broker   Broker       `mapstructure:"broker"`
}

type Broker struct {
	Brokers string `mapstructure:"brokers"`

	// Subscribers
	AC   cfg.Subscriber `mapstructure:"ac"`
	AECR cfg.Subscriber `mapstructure:"aecr"`
	AEVR cfg.Subscriber `mapstructure:"aevr"`
	APRR cfg.Subscriber `mapstructure:"aprr"`
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
	cfg.Client.Addr = cfg.Client.Host + ":" + strconv.Itoa(cfg.Client.Port)
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

	if env == "prod" {
		cfg.Client.Addr = "https://" + cfg.Client.Addr
	} else {
		cfg.Client.Addr = "http://" + cfg.Client.Addr
	}

	return &cfg, nil
}
