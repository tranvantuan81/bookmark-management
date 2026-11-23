package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranvantuan81/bookmark-management/internal/handler"
	"github.com/tranvantuan81/bookmark-management/internal/service"
)

// Engine defines the interface for the HTTP engine.
type Engine interface {
	// Start runs the HTTP engine on the configured application port.
	Start() error
	// ServeHTTP handles an HTTP request using the configured Gin engine.
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *Config
}

// NewEngine creates and configures an HTTP engine and its routes.
func NewEngine(cfg *Config) Engine {
	e := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	e.initRoutes()
	return e
}

// ServeHTTP handles an HTTP request using the configured Gin engine.
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// Start runs the HTTP engine on the configured application port.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.AppPort))
}

// initRoutes initializes the HTTP routes for the engine.
func (e *engine) initRoutes() {
	healthCheckSvc := service.NewHealthCheck(e.cfg.AppName, e.cfg.InstanceID)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
}
