package server

import "github.com/spf13/viper"

// Configuration contains the configuration to connect to the database.
type Configuration struct {
	AppName string

	Host string
	Port int

	ShutdownTimeout int
}

// NewConfig returns an instance to the configuration.
func NewConfig() *Configuration {
	return &Configuration{
		AppName: viper.GetString("APP_NAME"),

		Host: viper.GetString("HTTP_HOST"),
		Port: viper.GetInt("HTTP_PORT"),

		ShutdownTimeout: viper.GetInt("HTTP_TIMEOUT_SHUTDOWN"),
	}
}
