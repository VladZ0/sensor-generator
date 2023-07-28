package main

import (
	"context"
	"log"
	"sensors-generator/config"
	_ "sensors-generator/docs"
	"sensors-generator/internal/app"
	"sensors-generator/pkg/logging"
)

func main() {
	log.Print("Create configs")
	cfg := config.GetConfig()

	// I use this logger, because it is better for me.
	log.Print("Create logger.")
	logging.Init(cfg.AppConfig.LogLevel, false)
	logger := logging.GetLogger()

	logger.Info("Create app.")
	app, err := app.NewApp(context.Background(), cfg, logger)
	if err != nil {
		logger.Fatalf("Application failed, due to error: %v", err)
	}

	logger.Info("Start application.")
	app.Run()
}
