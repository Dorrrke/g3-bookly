package main

import (
	"github.com/Dorrrke/g3-bookly/internal/config"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	"github.com/Dorrrke/g3-bookly/internal/server"
)

func main() {
	cfg := config.ReadConfig()
	log := logger.Get(cfg.Debug)
	log.Debug().Any("cfg", cfg).Send()

	serv := server.New(*cfg)
	err := serv.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("server fatal error")
	}
	log.Info().Msg("server stoped")
}
