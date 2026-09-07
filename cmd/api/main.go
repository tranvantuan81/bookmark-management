package main

import (
	"github.com/tranvantuan81/bookmark-management/internal/api"
	"github.com/tranvantuan81/bookmark-management/internal/config"
)

func main() {
	// load config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	app := api.NewEngine(cfg)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
