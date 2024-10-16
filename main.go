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
	cfg := config.LoadAppConfig(yamlFile)
	serviceEnv := config.LoadЕnvConfig()

	// setup : service logger
	lgr := logger.Setup(serviceEnv)

	lgr.Info().
		Str("name", serviceName).
		Str("environment", serviceEnv.Name).
		Str("started", upTime).
		Str("version", version).
		Str("logLevel", serviceEnv.LogLevel).
		Str("port", serviceEnv.Port).
		Str("dialect", cfg.Database.Dialect).
		Str("Dbname", cfg.Database.Dbname).
		Str("Host", cfg.Database.Host).
		Str("Password", cfg.Database.Password).
		Str("Username", cfg.Database.Username).
		Msg("service details, starting the service")

	// setup : start service
	server.StartService(serviceEnv, cfg, lgr)

	lgr.Fatal().Msg("service stopped")
	return nil
}
