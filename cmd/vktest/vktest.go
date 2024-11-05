package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"

	"github.com/lesienchik/vk__test/internal/api"
	"github.com/lesienchik/vk__test/internal/config"
	"github.com/lesienchik/vk__test/internal/logic"
	"github.com/lesienchik/vk__test/internal/storage"
	postgres "github.com/lesienchik/vk__test/pkg/db"
	"github.com/lesienchik/vk__test/pkg/email"
)

// @title Vktest application
// @version 2.0
// @description The backend service for the site vktest.

// @host localhost:9100
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	data, err := os.ReadFile("./local_files/config.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := new(config.Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		log.Fatal(err)
	}

	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	logger := log.New()
	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger.SetLevel(logLevel)

	logger.Info("config initialization successfully")

	db, err := postgres.ConnectToDb(&cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("connection to DB successfully")

	storage := storage.New(logger, db)
	email := email.New(&cfg.Email)
	logic := logic.New(cfg.Logic.SecretKey, logger, email, storage)
	api := api.New(cfg, logger, logic)

	termChan, errChan := make(chan os.Signal, 1), make(chan error, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := api.Start(); err != nil {
			errChan <- err
		}
	}()
	logger.Info("api successfully started")
	logger.Info("vktest service has been successfully launched")

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-termChan:
		logger.Info("vktest service has been successfully stopped")
		api.Shutdown()
	}
}
