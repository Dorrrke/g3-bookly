package config

import (
	"flag"
	"fmt"
)

const (
	defaultAddr        = "localhost"
	defaultPort        = 8080
	defaultDBDsn       = "postgres://user:password@localhost:5432/course?sslmode=disable"
	defaultMigratePath = "migrations"
)

type Config struct {
	Addr        string
	Debug       bool
	DBDsn       string
	MigratePath string
}

func ReadConfig() *Config {
	var host, dbDsn, migratePath string
	var port int
	var debug bool
	flag.StringVar(&host, "addr", defaultAddr, "flag to set the server startup host")
	flag.IntVar(&port, "port", defaultPort, "flag to set the server startup port")
	flag.BoolVar(&debug, "debug", false, "flag to set Debug logger level")
	flag.StringVar(&dbDsn, "db", defaultDBDsn, "database connection addres")
	flag.StringVar(&migratePath, "m", defaultMigratePath, "path to migrations")
	flag.Parse()

	return &Config{
		Addr:        fmt.Sprintf("%s:%d", host, port),
		Debug:       debug,
		DBDsn:       dbDsn,
		MigratePath: migratePath,
	}
}
