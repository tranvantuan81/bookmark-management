package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranvantuan81/bookmark-management/internal/app/service"
)

// HealthCheck interface for health check handler
type HealthCheck interface {
	HealthCheck(c *gin.Context)
}

type healthCheckHandler struct {
	healthCheckService service.HealthCheck
}

// NewHealthCheck creates a new health check handler
func NewHealthCheck(healthCheckSvc service.HealthCheck) HealthCheck {
	return &healthCheckHandler{
		healthCheckService: healthCheckSvc,
	}
}

// HealthCheck Checks the health of the service
// @Summary      Check health
// @Description  Check health for app
// @Tags         Health check
// @Accept       application/json
// @Produce      application/json
// @Success      200  {object} 	map[string]string
// @Router       /health-check [GET]
func (s *healthCheckHandler) HealthCheck(c *gin.Context) {
	res, err := s.healthCheckService.HealthCheck(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, res)
}
