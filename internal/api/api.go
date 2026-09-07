package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranvantuan81/bookmark-management/internal/config"
	"github.com/tranvantuan81/bookmark-management/internal/handler"
	"github.com/tranvantuan81/bookmark-management/internal/service"
)

// Engine is the interface for the application
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *config.Config
}

// NewEngine creates a new engine
func NewEngine(cfg *config.Config) Engine {
	app := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	app.initRoutes()

	return app
}

// Start starts the application
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP implements the http.Handler interface
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// initRoutes initializes the routes
func (e *engine) initRoutes() {

	// genpass handler
	genPassSvc := service.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassSvc)

	healthCheckSev := service.NewHealthCheck(e.cfg)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSev)

	e.app.GET("/genpass", genPassHandler.GeneratePassword)
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
}
