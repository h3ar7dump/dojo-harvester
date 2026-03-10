package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/dojo-harvester/backend/internal/api"
	"github.com/dojo-harvester/backend/internal/config"
	"github.com/dojo-harvester/backend/internal/logger"
	"github.com/dojo-harvester/backend/internal/storage"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err) // Cannot use logger yet
	}

	// Initialize logger
	if err := logger.Init(&cfg.Logger); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Get()
	log.Info("Starting Dojo Harvester Agent Backend")

	// Initialize storage
	store, err := storage.NewStore(&cfg.Storage)
	if err != nil {
		log.Fatal("Failed to initialize storage", zap.Error(err))
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Error("Error closing storage", zap.Error(err))
		}
	}()

	// Initialize and start server
	server := api.NewServer(cfg, store)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}
