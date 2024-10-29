package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dorrrke/g3-bookly/internal/config"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	"github.com/Dorrrke/g3-bookly/internal/server"
	"github.com/Dorrrke/g3-bookly/internal/storage"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log := logger.Get(cfg.Debug)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		<-c

		log.Debug().Msg("ctx cancel; chatch os signal")
		cancel()
	}()

	log.Debug().Any("cfg", cfg).Send()
	var stor server.Storage

	if err := storage.Migrations(cfg.DBDsn, cfg.MigratePath); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}
	stor, err = storage.NewDB(context.TODO(), cfg.DBDsn)
	if err != nil {
		log.Error().Err(err).Msg("connecting to data base failed")
		stor = storage.New()
	}
	serv := server.New(*cfg, stor)
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return serv.Run()
	})

	group.Go(func() error {
		<-gCtx.Done()
		return serv.ShutdownServer()
	})

	if err := group.Wait(); err != nil {
		log.Info().Str("stoping reason", err.Error()).Msg("Server stoped")
		return
	}
	log.Info().Msg("server stoped")
}
