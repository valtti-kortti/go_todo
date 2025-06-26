package main

import (
	"go_todo/internal/api"
	"go_todo/internal/repo"
	"go_todo/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go_todo/internal/config"
	customLogger "go_todo/internal/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func main() {
	err := godotenv.Load(config.EnvConfigPath)
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Error loading config"))
	}

	logger, err := customLogger.NewLogger(cfg.Loglevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error initializing logger"))
	}

	repository := repo.NewRepository()

	serverInstance := service.NewService(repository, logger)

	app := api.NewRouters(&api.Routers{Service: serverInstance}, cfg.Rest.Token)

	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
	
}
