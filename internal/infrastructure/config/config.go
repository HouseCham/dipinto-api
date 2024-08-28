package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config represents the configuration of the application found within the config.json file
type Config struct {
	Server struct {
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	JWT struct {
		SecretKey string
	}
	Client struct {
		Origin string
	}
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		return nil, err
	}

	return &cfg, nil
}
