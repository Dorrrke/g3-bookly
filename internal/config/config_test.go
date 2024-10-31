package config

import (
	"flag"
	"os"
	"testing"

	"github.com/Dorrrke/g3-bookly/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	logger.Get(false)
	type want struct {
		cfg Config
	}
	type test struct {
		name  string
		flags []string
		env   func()
		want  want
	}
	tests := []test{
		{
			name: "succssesful flags call",
			flags: []string{
				"test", "-addr", "134.222.12.12",
				"-port", "1234", "-debug", "-db", "testDBURL",
				"-m", "/test/migrate/path",
			},
			want: want{
				cfg: Config{
					Addr:        "134.222.12.12:1234",
					Debug:       true,
					DBDsn:       "testDBURL",
					MigratePath: "/test/migrate/path",
				},
			},
		},
		{
			name:  "succssesful env call",
			flags: []string{"test"},
			env: func() {
				t.Setenv("SERVER_HOST", "111.11.11.12")
				t.Setenv("SERVER_PORT", "9695")
				t.Setenv("DB_DSN", "test://db:dsn")
				t.Setenv("MIGRATE_PATH", "/test/migrate/path")
			},
			want: want{
				cfg: Config{
					Addr:        "111.11.11.12:9695",
					Debug:       false,
					DBDsn:       "test://db:dsn",
					MigratePath: "/test/migrate/path",
				},
			},
		},
		{
			name:  "default call",
			flags: []string{"test"},
			want: want{
				cfg: Config{
					Addr:        "localhost:8080",
					Debug:       false,
					DBDsn:       "postgres://user:password@localhost:5432/course?sslmode=disable",
					MigratePath: "migrations",
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tc.flags
			if tc.env != nil {
				tc.env()
				defer os.Unsetenv("SERVER_HOST")
				defer os.Unsetenv("SERVER_PORT")
				defer os.Unsetenv("DB_DSN")
				defer os.Unsetenv("MIGRATE_PATH")
			}
			cfg, err := ReadConfig()
			assert.NoError(t, err)
			assert.Equal(t, tc.want.cfg, *cfg)
		})
	}
}
