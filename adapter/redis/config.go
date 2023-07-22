package redis

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host       string
	Password   string
	Port       int
	Index      int
	Expiration time.Duration // in seconds
}

func NewConfig() *Config {
	return &Config{
		Host:       viper.GetString("REDIS_HOST"),
		Password:   viper.GetString("REDIS_PASSWORD"),
		Port:       viper.GetInt("REDIS_PORT"),
		Index:      viper.GetInt("REDIS_INDEX"),
		Expiration: viper.GetDuration("REDIS_EXPIRATION"),
	}
}
