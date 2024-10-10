package main

import (
	"github.com/Dorrrke/g3-bookly/internal/server"
)

func main() {
	serv := server.New(":8080")
	err := serv.Run()
	if err != nil {
		panic(err)
		//log.Fatal().Err(err).Msg("server fatal error")
	}
	//log.Info().Msg("server stoped")
}
