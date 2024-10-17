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
		From       string   `yaml:"from"`
		Password   string   `yaml:"password"`
		Recipients []string `yaml:"recipients"`
		SmtpHost   string   `yaml:"smtp_host"`
		SmtpPort   string   `yaml:"smtp_port"`
	} `yaml:"email"`
	Slack struct {
		WebhookUrl string `yaml:"webhook_url"`
	} `yaml:"slack"`
	Database struct {
		Dialect  string `yaml:"dialect"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Dbname   string `yaml:"dbname"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

type ServiceEnv struct {
	Name     string // name of environment where this service is running
	Port     string // port on which this service runs, defaults to DefaultPort
	LogLevel string // logger level for the service
}

// LoadAppConfig reads the settings written to the yml file.
func LoadAppConfig(yamlFile embed.FS) *Config {
	var env *string
	if value := os.Getenv("WEB_APP_ENV"); value != "" {
		env = &value
	} else {
		env = flag.String("env", "develop", "To switch configurations.")
		flag.Parse()
	}

	file, fileErr := yamlFile.ReadFile(fmt.Sprintf(AppConfigPath, *env))
	if fileErr != nil {
		os.Exit(ErrExitStatus)
	}

	config := &Config{}
	if err := yaml.Unmarshal(file, config); err != nil {
		os.Exit(ErrExitStatus)
	}

	return config
}

// Load the service environment variable related to the service configuration.
func LoadEnvConfig() ServiceEnv {
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
