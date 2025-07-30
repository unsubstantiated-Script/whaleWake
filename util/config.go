package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DB_DRIVER               string        `mapstructure:"DB_DRIVER"`
	DB_SOURCE               string        `mapstructure:"DB_SOURCE"`
	SERVER_ADDRESS          string        `mapstructure:"SERVER_ADDRESS"`
	PASETO_SYMMETRIC_KEY    string        `mapstructure:"PASETO_SYMMETRIC_KEY"`
	ACCESS_TOKEN_EXPIRATION time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRATION"`
}

func LoadConfig(path string) (config Config, err error) {
	// Load environment variables from the specified path
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
