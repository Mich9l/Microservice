package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"myservice/internal/config"
	"myservice/internal/infrastructure/database"
	"myservice/internal/infrastructure/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "критическая ошибка запуска: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Инициализация логгера
	log, err := logger.New()
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	defer func() { _ = log.Sync() }()

	// 2. Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// 3. Подключение к PostgreSQL
	db, err := database.NewPostgres(cfg)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer db.Close()

	log.Info("успешное подключение к PostgreSQL",
		zap.String("host", cfg.DBHost),
		zap.String("port", cfg.DBPort),
		zap.String("database", cfg.DBName),
	)

	return nil
}
