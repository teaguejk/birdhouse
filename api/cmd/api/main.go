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
	"api/internal/command"
	"api/internal/device"
	"api/internal/server"
	"api/pkg/ai"
	"api/pkg/database"
	"api/pkg/logging"
	"api/pkg/mqtt"
	"api/pkg/oauth"
	"api/pkg/storage"
)

func main() {
	ctx := context.Background()

	logger := logging.New()
	ctx = logging.WithLogger(ctx, logger)

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.dev.json"
	}

	cfg, err := api.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewDatabase(ctx, cfg.Database, logger.WithField("component", "database"))
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	storage, err := storage.NewProvider(ctx, cfg.Storage)
	if err != nil {
		log.Fatalf("failed to initialize storage provider: %v", err)
	}

	aiClient, err := ai.NewClient(cfg.AI)
	if err != nil {
		log.Fatalf("failed to initialize ai client: %v", err)
	}
	if !aiClient.IsConfigured() {
		logger.Warn("ai was not configured, ai features will be unavailable")
	}

	oauthVerifier, err := oauth.NewVerifier(cfg.OAuth)
	if err != nil {
		log.Fatalf("failed to initialize oauth verifier: %v", err)
	}

	publisher, mqttClient := mqtt.NewFromConfig(cfg.MQTT, logger.WithField("component", "mqtt"))

	sEnv := &api.ServerEnv{
		Logger:        logger.WithField("component", "server"),
		Config:        cfg.Server,
		OAuthVerifier: oauthVerifier,
	}

	repos := api.InitRepositories(db)
	services := api.InitServices(repos, logger, storage, aiClient, publisher)
	handlers := api.InitHandlers(services, logger, db)

	if mqttClient != nil {
		ackSub := command.NewAckSubscriber(logger.WithField("subscriber", "ack"), repos.Command)
		if err := ackSub.Subscribe(mqttClient); err != nil {
			log.Fatalf("failed to subscribe to ack topic: %v", err)
		}

		statusSub := device.NewStatusSubscriber(logger.WithField("subscriber", "status"), repos.Device)
		if err := statusSub.Subscribe(mqttClient); err != nil {
			log.Fatalf("failed to subscribe to status topic: %v", err)
		}
	}

	srv := server.New(sEnv, services)
	srv.RegisterHandler(handlers.Auth)
	srv.RegisterHandler(handlers.Command)
	srv.RegisterHandler(handlers.Device)
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

	if mqttClient != nil {
		mqttClient.Disconnect()
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server failed to shutdown: %v", err)
	}

	db.Close()

	logger.Info("successfully shutdown")
}
