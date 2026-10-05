package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tranvantuan81/bookmark-management/docs"
	"github.com/tranvantuan81/bookmark-management/internal/config"
	"github.com/tranvantuan81/bookmark-management/internal/handler"
	"github.com/tranvantuan81/bookmark-management/internal/repository"
	"github.com/tranvantuan81/bookmark-management/internal/service"
)

// Engine is the interface for the application
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app         *gin.Engine
	cfg         *config.Config
	redisClient *redis.Client
}

// NewEngine creates a new engine
func NewEngine(cfg *config.Config, redisClient *redis.Client) Engine {
	app := &engine{
		app:         gin.Default(),
		cfg:         cfg,
		redisClient: redisClient,
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

	healthCheckRepo := repository.NewPing(e.redisClient)
	healthCheckSev := service.NewHealthCheck(e.cfg, healthCheckRepo)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSev)

	urlStorage := repository.NewURLStorage(e.redisClient)
	urlService := service.NewShortenUrl(urlStorage)
	urlHandler := handler.NewShortenLink(urlService)

	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
	e.app.POST("/v1/links/shorten", urlHandler.CreateShortenLink)
	e.app.GET("/v1/links/redirect/:code", urlHandler.Redirect)
}
