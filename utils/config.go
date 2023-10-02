package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBdriver       string        `mapstructure:"DB_DRIVER"`
	DB_Source      string        `mapstructure:"DB_SOURCE"`
	DB_Source_Live string        `mapstructure:"DB_SOURCE_LIVE"`
	ServerPort     string        `mapstructure:"SERVER_PORT"`
	SymetricKey    string        `mapstructure:"SYMETRIC_KEY"`
	ExpDuration    time.Duration `mapstructure:"EXP_DURATION"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
