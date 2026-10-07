package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tranvantuan81/bookmark-management/docs"
	_ "github.com/tranvantuan81/bookmark-management/docs"
	handler2 "github.com/tranvantuan81/bookmark-management/internal/app/handler"
	repository2 "github.com/tranvantuan81/bookmark-management/internal/app/repository"
	service2 "github.com/tranvantuan81/bookmark-management/internal/app/service"
	"github.com/tranvantuan81/bookmark-management/internal/config"
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

	healthCheckRepo := repository2.NewPing(e.redisClient)
	healthCheckSev := service2.NewHealthCheck(e.cfg, healthCheckRepo)
	healthCheckHandler := handler2.NewHealthCheck(healthCheckSev)

	urlStorage := repository2.NewURLStorage(e.redisClient)
	urlService := service2.NewShortenUrl(urlStorage, service2.NewCodeGen())
	urlHandler := handler2.NewShortenLink(urlService)

	// Init swagger routes
	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
	e.app.POST("/v1/links/shorten", urlHandler.CreateShortenLink)
	e.app.GET("/v1/links/redirect/:code", urlHandler.Redirect)
}
