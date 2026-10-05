package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	log.Debug().Msg("Hello World debug")
	log.Info().Msg("Hello World info")
	log.Error().Msg("Hello World error")
}
