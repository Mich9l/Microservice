package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"myservice/internal/config"
	"myservice/internal/infrastructure/database"
	"myservice/internal/infrastructure/logger"
	"myservice/internal/infrastructure/repository"
	grpcserver "myservice/internal/infrastructure/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "критическая ошибка запуска: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	log, err := logger.New()
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	defer func() { _ = log.Sync() }()

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

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

	repo := repository.NewRecipeRepository(db)
	grpcAddr := ":" + cfg.GRPCPort

	grpcSrv, err := grpcserver.New(grpcAddr, repo, log)
	if err != nil {
		return fmt.Errorf("gRPC server: %w", err)
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Info("gRPC server starting", zap.String("addr", grpcAddr))
		serveErr <- grpcSrv.Start()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Info("shutdown signal received", zap.String("signal", sig.String()))
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("gRPC serve: %w", err)
		}
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := grpcSrv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("gRPC shutdown: %w", err)
	}

	log.Info("gRPC server stopped")
	return nil
}