package config

type ContextKey string

const RequestIdentifier = "X-Request-ID"

const (
	// AppConfigPath is the path of application.yml.
	AppConfigPath = "resources/config/application.%s.yml"
)
