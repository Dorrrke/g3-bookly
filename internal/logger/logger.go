package logger

import (
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

var once sync.Once //nolin:gochecknoglobals //singletone

var log zerolog.Logger //nolin:gochecknoglobals //singletone

func Get(flags ...bool) zerolog.Logger {
	once.Do(func() {
		zerolog.TimestampFieldName = "Time"
		zerolog.LevelFieldName = "Level"
		zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
			return file + ":" + strconv.Itoa(line)
		}
		if flags[0] {
			log = zerolog.New(os.Stdout).
				Level(zerolog.DebugLevel).
				With().
				Timestamp().
				Caller().
				Logger().
				Output(zerolog.ConsoleWriter{Out: os.Stderr})
		} else {
			log = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
		}
	})
	return log
}
