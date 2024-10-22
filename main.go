package main

import (
	"embed"
	"time"

	"github.com/plyovchev/notifications-service/internal/config"
	"github.com/plyovchev/notifications-service/internal/logger"
	"github.com/plyovchev/notifications-service/internal/server"
)

//go:embed resources/config/application.*.yml
var yamlFile embed.FS

const (
	serviceName = "notification-gateway"
)

// Passed while building from  make file.
var version string

func main() {
	upTime := time.Now().UTC().Format(time.RFC3339)
	// setup : read environmental configurations
	// setup : service logger
	cfg := config.LoadAppConfig(yamlFile)
	serviceEnv := config.LoadEnvConfig()

	// setup : service logger
	lgr := logger.Setup(serviceEnv)

	lgr.Info().
		Str("name", serviceName).
		Str("environment", serviceEnv.Name).
		Str("started", upTime).
		Str("version", version).
		Str("logLevel", serviceEnv.LogLevel).
		Str("port", serviceEnv.Port).
		Msg("service details, starting the service")

	// setup : start service
	server.StartService(serviceEnv, cfg, lgr)

	lgr.Fatal().Msg("service stopped")
}
