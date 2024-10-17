package main

import (
	"context"

	"github.com/Dorrrke/g3-bookly/internal/config"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	"github.com/Dorrrke/g3-bookly/internal/server"
	"github.com/Dorrrke/g3-bookly/internal/storage"
)

func main() {
	cfg := config.ReadConfig()
	log := logger.Get(cfg.Debug)
	log.Debug().Any("cfg", cfg).Send()
	var stor server.Storage

	if err := storage.Migrations(cfg.DBDsn, cfg.MigratePath); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}
	stor, err := storage.NewDB(context.TODO(), cfg.DBDsn)
	if err != nil {
		log.Error().Err(err).Msg("connecting to data base failed")
		stor = storage.New()
	}

	serv := server.New(*cfg, stor)
	err = serv.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("server fatal error")
	}
	log.Info().Msg("server stoped")
}
