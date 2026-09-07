package config

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config struct
type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark-management" envconfig:"SERVICE_NAME"`
	InstanceID  string `default:"" envconfig:"INSTANCE_ID"`
}

// NewConfig creates a new config
func NewConfig() (*Config, error) {
	config := &Config{}

	if config.InstanceID == "" {
		config.InstanceID = uuid.NewString()
	}

	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}
	return config, err
}
