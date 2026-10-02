package main

import (
	"github.com/tranvantuan81/bookmark-management/internal/api"
	"github.com/tranvantuan81/bookmark-management/internal/config"
	"github.com/tranvantuan81/bookmark-management/pkg/redis"
)

// @title           Bookmark Management API
// @version         1.0
// @description     API for managing bookmark.
// @BasePath /
func main() {
	// load config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	defer redisClient.Close()

	app := api.NewEngine(cfg, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
