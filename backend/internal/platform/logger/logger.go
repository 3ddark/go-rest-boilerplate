package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// L, projenin her yerinden erişilebilecek global logger'dır.
var L zerolog.Logger

func Init() {
	// Geliştirme ortamı için daha okunaklı, renkli loglar
	// Production'da ise JSON formatında loglama yapmak daha iyidir.
	// Bunu bir ortam değişkeniyle kontrol edebiliriz.
	isDev := os.Getenv("APP_ENV") == "development"

	if isDev {
		L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		L = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
