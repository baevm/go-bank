package config

import "github.com/spf13/viper"

type Config struct {
	DB_DSN   string `mapstructure:"DB_DSN"`
	SRV_ADDR string `mapstructure:"SRV_ADDR"`
}

func Load(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)

	return &cfg, err
}
