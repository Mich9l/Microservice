package config

import (
	"fmt"
	"os"
	"strings"
)

func loadDotEnv(paths ...string) {
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for line := range strings.SplitSeq(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}
		return
	}
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	GRPCPort   string
}

func LoadConfig() (Config, error) {
	loadDotEnv(".env")

	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		GRPCPort:   os.Getenv("GRPC_PORT"),
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

	if cfg.GRPCPort == "" {
		cfg.GRPCPort = "50051"
	}

	return cfg, nil
}
