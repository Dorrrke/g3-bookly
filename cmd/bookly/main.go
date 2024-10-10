package main

import (
	"g3-bookly/internal/config"
	"g3-bookly/internal/logger"
	"g3-bookly/internal/server"
	"g3-bookly/internal/storage"
)

func main() {
	cfg := config.ReadConfig()
	log := logger.Get(cfg.Debug)
	log.Debug().Any("cfg", cfg).Send()

	stor := storage.New()

	serv := server.New(*cfg, stor)
	err := serv.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("server fatal error")
	}
	log.Info().Msg("server stoped")
}
