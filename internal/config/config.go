package config

import (
	"cmp"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	defaultAddr        = "localhost"
	defaultPort        = 8080
	defaultDBDsn       = "postgres://user:password@localhost:5432/course?sslmode=disable"
	defaultMigratePath = "migrations"
	defaultAuthHost    = "localhost:9090"
)

type Config struct {
	Addr        string
	Debug       bool
	DBDsn       string
	MigratePath string
	AuthHost    string
}

func ReadConfig() (*Config, error) {
	var host, dbDsn, migratePath, authHost string
	var port int
	var debug bool
	flag.StringVar(&host, "addr", defaultAddr, "flag to set the server startup host")
	flag.StringVar(&authHost, "g", defaultAuthHost, "flag to set the auth service host")
	flag.IntVar(&port, "port", defaultPort, "flag to set the server startup port")
	flag.BoolVar(&debug, "debug", false, "flag to set Debug logger level")
	flag.StringVar(&dbDsn, "db", defaultDBDsn, "database connection addres")
	flag.StringVar(&migratePath, "m", defaultMigratePath, "path to migrations")
	flag.Parse()

	host = cmp.Or(os.Getenv("SERVER_HOST"), host)
	p := cmp.Or(os.Getenv("SERVER_PORT"), strconv.Itoa(port))
	port, err := strconv.Atoi(p)
	if err != nil {
		return nil, err
	}
	dbDsn = cmp.Or(os.Getenv("DB_DSN"), dbDsn)
	authHost = cmp.Or(os.Getenv("AUTH_DSN"), authHost)
	migratePath = cmp.Or(os.Getenv("MIGRATE_PATH"), migratePath)
	return &Config{
		Addr:        fmt.Sprintf("%s:%d", host, port),
		Debug:       debug,
		DBDsn:       dbDsn,
		MigratePath: migratePath,
		AuthHost:    authHost,
	}, nil
}
