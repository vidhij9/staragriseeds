package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all the configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
	SMTP     SMTPConfig
}

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Port int
	Host string
}

// DatabaseConfig holds the configuration for the database
type DatabaseConfig struct {
	TablePrefix string
}

// AWSConfig holds the configuration for AWS services
type AWSConfig struct {
	Region string
}

// SMTPConfig holds the configuration for SMTP
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Load reads the configuration from a file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	viper.AddConfigPath("./config") // look for config in the config directory

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("No config file found. Using defaults and environment variables.")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config

	// Server configuration
	config.Server.Port = viper.GetInt("server.port")
	config.Server.Host = viper.GetString("server.host")

	// Database configuration
	config.Database.TablePrefix = viper.GetString("database.tablePrefix")

	// AWS configuration
	config.AWS.Region = viper.GetString("aws.region")

	// SMTP configuration
	config.SMTP.Host = viper.GetString(("smtp.username"))
	config.SMTP.Port = viper.GetInt(("smtp.port"))
	config.SMTP.Username = viper.GetString(("smtp.username"))
	config.SMTP.Password = viper.GetString(("smtp.password"))

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig checks if the required configuration fields are set
func validateConfig(config *Config) error {
	if config.Server.Port == 0 {
		return fmt.Errorf("server port is required")
	}
	if config.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if config.AWS.Region == "" {
		return fmt.Errorf("AWS region is required")
	}
	if config.SMTP.Host == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if config.SMTP.Port == 0 {
		return fmt.Errorf("SMTP port is required")
	}
	if config.SMTP.Username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if config.SMTP.Password == "" {
		return fmt.Errorf("SMTP password is required")
	}
	return nil
}
