package main

import "github.com/tranvantuan81/bookmark-management/internal/api"

func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg)
	err = app.Start()
	if err != nil {
		panic(err)
	}

}
