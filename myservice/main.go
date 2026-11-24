package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"myservice/internal/application/processor"
	"myservice/internal/infrastructure/kafka"
	"myservice/internal/infrastructure/logger"
	"myservice/internal/infrastructure/repository"
	grpcserver "myservice/internal/infrastructure/server"
	httpserver "myservice/internal/server"
)

type config struct {
	DatabaseURL  string
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
	GRPCAddr     string
	HTTPAddr     string
}

func loadConfig() (config, error) {
	cfg := config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
		GRPCAddr:     os.Getenv("GRPC_ADDR"),
		HTTPAddr:     os.Getenv("HTTP_ADDR"),
	}
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		cfg.KafkaBrokers = []string{v}
	}
	return cfg, nil
}

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
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. Подключение к PostgreSQL
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer pool.Close()

	// 4. Создание Repository
	recipeRepo := repository.NewRecipeRepository(pool)

	// 5. Создание Recipe Processor
	recipeProcessor := processor.NewRecipeProcessor(recipeRepo, log)

	// 6. Инициализация Kafka Consumer
	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID, recipeProcessor, log)
	if err != nil {
		return fmt.Errorf("kafka consumer: %w", err)
	}

	// 7. Инициализация gRPC Server
	grpcSrv, err := grpcserver.New(cfg.GRPCAddr, recipeRepo, log)
	if err != nil {
		return fmt.Errorf("grpc server: %w", err)
	}

	// 8. Инициализация HTTP Server
	httpSrv, err := httpserver.New(cfg.HTTPAddr, recipeProcessor, log)
	if err != nil {
		return fmt.Errorf("http server: %w", err)
	}

	// 9. Запуск серверов
	errCh := make(chan error, 3)

	go func() {
		errCh <- consumer.Run(ctx)
	}()
	go func() {
		errCh <- grpcSrv.Start()
	}()
	go func() {
		errCh <- httpSrv.Start()
	}()

	// 10. Ожидание сигнала завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil {
			log.Error("компонент завершился с ошибкой", zap.Error(err))
		}
	case sig := <-sigCh:
		log.Info("получен сигнал завершения", zap.String("signal", sig.String()))
	}

	// 11. Graceful Shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	cancel()

	consumer.Shutdown()
	if err := grpcSrv.Shutdown(shutdownCtx); err != nil {
		log.Error("grpc shutdown", zap.Error(err))
	}
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown", zap.Error(err))
	}

	return nil
}