package main

import (
	"embed"
	"os"
	"time"

	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/server"
)

//go:embed resources/config/application.*.yml
var yamlFile embed.FS

const (
	serviceName = "notification-gateway"
	defaultPort = "8080"
)

// Passed while building from  make file.
var version string

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	upTime := time.Now().UTC().Format(time.RFC3339)
	// setup : read environmental configurations
	// setup : service logger
	config, _ := config.LoadAppConfig(yamlFile)

	// setup : service logger
	lgr := logger.Setup(config)

	lgr.Info().Msgf("config email %s", config)

	lgr.Info().
		Str("name", serviceName).
		Str("environment", config.Name).
		Str("started", upTime).
		Str("version", version).
		Str("logLevel", config.LogLevel).
		Msg("service details, starting the service")

	// setup : start service
	server.StartService(config, lgr)

	lgr.Fatal().Msg("service stopped")
	return nil
}
