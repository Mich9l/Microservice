package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	var missing []string
	if cfg.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.DBPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if cfg.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.DBName == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return Config{}, fmt.Errorf("required environment variables are not set: %v", missing)
	}

	return cfg, nil
}
