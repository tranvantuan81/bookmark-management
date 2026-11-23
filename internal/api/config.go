package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppName    string `default:"bookmark_service" envconfig:"APP_NAME"`
	AppPort    int    `default:"8080" envconfig:"APP_PORT"`
	InstanceID string `default:"" envconfig:"INSTANCE_ID"`
}

// NewConfig loads application configuration from environment variables.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg, nil
}
