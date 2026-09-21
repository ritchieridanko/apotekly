package configs

type Storage struct {
	Provider string `mapstructure:"provider"`
	Cloud    string `mapstructure:"cloud"`
	Key      string `mapstructure:"key"`
	Secret   string `mapstructure:"secret"`
}
