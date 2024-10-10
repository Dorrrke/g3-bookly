package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Addr  string
	Debug bool
}

func ReadConfig() *Config {
	var host, port string
	var debug bool
	flag.StringVar(&host, "addr", "localhost", "flag to set the server startup host")
	flag.StringVar(&port, "port", "8080", "flag to set the server startup port")
	flag.BoolVar(&debug, "debug", false, "flag to set Debug logger level")
	flag.Parse()

	return &Config{
		Addr:  fmt.Sprintf("%s:%s", host, port),
		Debug: debug,
	}
}
