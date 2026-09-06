package service

import "github.com/tranvantuan81/bookmark-management/internal/config"

// HealthCheck interface for health check service
//go:generate mockery --name=HealthCheck --filename=health_check.go
type HealthCheck interface {
	CheckHealth() (Response, error)
}

type healthCheckService struct {
	cfg *config.Config
}

// NewHealthCheck creates a new health check service
func NewHealthCheck(cfg *config.Config) HealthCheck {
	return &healthCheckService{
		cfg: cfg,
	}
}

// Response struct for health check response
type Response struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// CheckHealth checks the health of the service
func (s *healthCheckService) CheckHealth() (Response, error) {
	// load config
	//cfg, _ := config.NewConfig()

	return Response{
		Message:     "OK",
		ServiceName: s.cfg.ServiceName,
		InstanceID:  s.cfg.InstanceID,
	}, nil

	//return res, nil
}
