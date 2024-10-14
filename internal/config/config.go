package config

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ErrExitStatus represents the error status in this application.
const ErrExitStatus int = 2

// Config represents the composition of yml settings.
type Config struct {
	LogLevel string `yaml:"log_level"`
	Name     string
	Port     string
	Email    struct {
		From       string
		Password   string
		Recipients []string
		SmtpHost   string `yaml:"smtp_host"`
		SmtpPort   string `yaml:"smtp_port"`
	}
	Slack struct {
		WebhookUrl string `yaml:"webhook_url"`
	}
	Kafka struct {
	}
}

// LoadAppConfig reads the settings written to the yml file
func LoadAppConfig(yamlFile embed.FS) (*Config, string) {
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

	return config, *env
}
