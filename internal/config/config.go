package config

import (
	"flag"
	"fmt"
)

type Cfg struct {
	Addr  string
	Debug bool
}

func ReadConfig() *Cfg {
	var host, port string
	var debug bool
	flag.StringVar(&host, "serv-addr", "localhost", "flag to set the server startup host")
	flag.StringVar(&port, "prt", "8080", "flag to set the server startup port")
	flag.BoolVar(&debug, "debug", false, "flag to set Debug logger level")
	flag.Parse()

	return &Cfg{
		Addr:  fmt.Sprintf("%s:%s", host, port),
		Debug: debug,
	}
}
