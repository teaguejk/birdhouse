package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api/internal/api"
	"api/internal/server"
	"api/pkg/ai"
	"api/pkg/database"
	"api/pkg/logging"
	"api/pkg/storage"
)

func main() {
	ctx := context.Background()

	logger := logging.New()
	ctx = logging.WithLogger(ctx, logger)

	cfg := api.NewDefaultConfig()

	// db, err := database.NewPostgresDB(ctx, dbCfg)
	db, err := database.NewPostgresDBWithName(ctx, cfg.Database, "birdhouse")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	storage, err := storage.NewProvider(ctx, cfg.Storage)
	if err != nil {
		log.Fatalf("failed to initialize storage provider: %v", err)
	}

	ai := ai.NewClient(cfg.AI)
	if !ai.IsConfigured() {
		logger.Warn("ai was not configured, ai features will be unavailable")
	}

	sEnv := &api.ServerEnv{
		Logger: logger,
		Config: cfg.Server,
	}

	repos := api.InitRepositories(db)
	services := api.InitServices(repos, logger, storage, ai)
	handlers := api.InitHandlers(services, logger, db)

	srv := server.New(sEnv, services)
	srv.RegisterHandler(handlers.Health)
	srv.RegisterHandler(handlers.Upload)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("starting shutdown...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server failed to shutdown: %v", err)
	}

	logger.Info("successfully shutdown")
}
