package client

import "github.com/spf13/viper"

// Configuration contains the configuration.
type Configuration struct {
	Timeout int // Timeout specifies a time limit for requests made by this Client.
}

// NewConfig returns an instance to the configuration.
func NewConfig() *Configuration {
	viper.SetDefault("HTTP_CLIENT_TIMEOUT", 10) // 10 seconds

	return &Configuration{
		Timeout: viper.GetInt("HTTP_CLIENT_TIMEOUT"),
	}
}
