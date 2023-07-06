package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DSN                string        `mapstructure:"DB_DSN"`
	SRV_ADDR              string        `mapstructure:"SRV_ADDR"`
	TOKEN_SYMMETRIC_KEY   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func Load(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg, nil
}
