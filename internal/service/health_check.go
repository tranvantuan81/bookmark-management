package service

import (
	"context"

	"github.com/tranvantuan81/bookmark-management/internal/config"
	"github.com/tranvantuan81/bookmark-management/internal/repository"
)

// CheckHealth interface for health check service
//
//go:generate mockery --name=HealthCheck --filename=health_check.go
type HealthCheck interface {
	HealthCheck(ctx context.Context) (Response, error)
}

type healthCheckService struct {
	cfg  *config.Config
	ping repository.Ping
}

// NewHealthCheck creates a new health check service
func NewHealthCheck(cfg *config.Config, ping repository.Ping) HealthCheck {
	return &healthCheckService{
		cfg:  cfg,
		ping: ping,
	}
}

// Response struct for health check response
type Response struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// HealthCheck checks the health of the service
func (s *healthCheckService) HealthCheck(ctx context.Context) (Response, error) {
	// load config
	//cfg, _ := config.NewConfig()

	err := s.ping.Ping(ctx)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Message:     "OK",
		ServiceName: s.cfg.ServiceName,
		InstanceID:  s.cfg.InstanceID,
	}, nil

	//return res, nil
}
