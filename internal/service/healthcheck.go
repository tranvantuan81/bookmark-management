package service

//go:generate mockery --name=HealthCheck --filename healthcheck.go --outpkg healthcheck_svc_mocks
type HealthCheck interface {
	// Check returns the current health status of the service.
	Check() (Response, error)
}

type healthCheck struct {
	serviceName string
	instanceID  string
}

type Response struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// NewHealthCheck creates a health check service for the specified service and instance.
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheck{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

// Check returns an OK response containing the service and instance identifiers.
func (h *healthCheck) Check() (Response, error) {
	return Response{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
