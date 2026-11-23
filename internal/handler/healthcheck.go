package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranvantuan81/bookmark-management/internal/service"
)

// HealthCheck defines the interface for the health check HTTP handler.
type HealthCheck interface {
	// HealthCheck handles health check HTTP requests.
	HealthCheck(c *gin.Context)
}

type healthCheck struct {
	svc service.HealthCheck
}

// NewHealthCheck creates a health check HTTP handler.
func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthCheck{
		svc: svc,
	}
}

// HealthCheck responds with the service health status.
func (h *healthCheck) HealthCheck(c *gin.Context) {
	result, err := h.svc.Check()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, service.Response{
		Message:     result.Message,
		ServiceName: result.ServiceName,
		InstanceID:  result.InstanceID,
	})
}
