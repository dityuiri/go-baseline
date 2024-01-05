package db

import (
	"github.com/spf13/viper"
)

type MigrationConfiguration struct {
	User     string
	Password string
}

// Configuration contains the configuration to connect to the database.
type Configuration struct {
	Host string
	Port int

	User     string
	Password string

	Name    string // The name of the database to connect to
	Schema  string
	Driver  string
	SSLMode string // Whether or not to use SSL (disable, require, verify-ca, verify-full)

	Migration *MigrationConfiguration

	MaxOpenConns    int // Maximum number of open connections to the database
	MaxIdleConns    int // Maximum number of connections in the idle connection pool
	ConnMaxLifetime int // Maximum amount of time a connection may be reused (in seconds)
	ConnMaxIdleTime int // Maximum amount of time a connection may be idle before being closed (in seconds)

	Echo bool // Enable echo to output SQL statements
}

// NewConfig returns an instance to the configuration.
func NewConfig() *Configuration {
	// https://www.alexedwards.net/blog/configuring-sqldb
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 25)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 5*60)  // 5 minutes
	viper.SetDefault("DB_CONN_MAX_IDLE_TIME", 1*60) // 1 minute

	viper.SetDefault("DB_ECHO", false)

	return &Configuration{
		Host: viper.GetString("DB_HOST"),
		Port: viper.GetInt("DB_PORT"),

		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),

		Name:    viper.GetString("DB_NAME"),
		Schema:  viper.GetString("DB_SCHEMA"),
		Driver:  viper.GetString("DB_DRIVER"),
		SSLMode: viper.GetString("DB_SSL_MODE"),

		Migration: &MigrationConfiguration{
			User:     viper.GetString("DB_MIGRATION_USER"),
			Password: viper.GetString("DB_MIGRATION_PASSWORD"),
		},

		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: viper.GetInt("DB_CONN_MAX_LIFETIME"),
		ConnMaxIdleTime: viper.GetInt("DB_CONN_MAX_IDLE_TIME"),

		Echo: viper.GetBool("DB_ECHO"),
	}
}
