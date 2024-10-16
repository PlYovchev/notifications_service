package config

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ErrExitStatus represents the error status in this application.
const (
	defaultPort       = "8080"
	ErrExitStatus int = 2
)

// Config represents the composition of yml settings.
type Config struct {
	Email struct {
		From       string
		Password   string
		Recipients []string
		SmtpHost   string `yaml:"smtp_host"`
		SmtpPort   string `yaml:"smtp_port"`
	}
	Slack struct {
		WebhookUrl string `yaml:"webhook_url"`
	}
	Database struct {
		Dialect  string
		Host     string
		Port     string
		Username string
		Dbname   string
		Password string
	}
}

type ServiceEnv struct {
	Name     string // name of environment where this service is running
	Port     string // port on which this service runs, defaults to DefaultPort
	LogLevel string // logger level for the service
}

// LoadAppConfig reads the settings written to the yml file
func LoadAppConfig(yamlFile embed.FS) *Config {
	var env *string
	if value := os.Getenv("WEB_APP_ENV"); value != "" {
		env = &value
	} else {
		env = flag.String("env", "develop", "To switch configurations.")
		flag.Parse()
	}

	file, err := yamlFile.ReadFile(fmt.Sprintf(AppConfigPath, *env))
	if err != nil {
		fmt.Printf("Failed to read application.%s.yml: %s", *env, err)
		os.Exit(ErrExitStatus)
	}

	config := &Config{}
	if err := yaml.Unmarshal(file, config); err != nil {
		fmt.Printf("Failed to read application.%s.yml: %s", *env, err)
		os.Exit(ErrExitStatus)
	}

	return config
}

// Load the service environment variable related to the service configuration
func LoadЕnvConfig() ServiceEnv {
	envName := os.Getenv("environment")
	if envName == "" {
		envName = "local"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logLevel := os.Getenv("logLevel")
	if logLevel == "" {
		logLevel = "info"
	}

	envConfigurations := ServiceEnv{
		Name:     envName,
		Port:     port,
		LogLevel: logLevel,
	}

	return envConfigurations
}
