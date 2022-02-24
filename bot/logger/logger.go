package logger

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog"
)

// Можно добавить в слой репозитория для дебажинга
type Logger struct {
	log zerolog.Logger
}

func (l *Logger) Error(err error, msg string) {
	l.log.Error().Err(err).Msg(msg)
}

func (l *Logger) Info(message string) {
	l.log.Info().Msg(message)
}

func (l *Logger) Debug(message string) {
	l.log.Debug().Msg(message)
}

// NewLogger returns instance of Logger, which satisfies domain.Logger interface
func NewLogger() *Logger {
	var log zerolog.Logger
	tempFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		log.Error().Err(err).Msg("there was an error creating a temporary file for our log")
	}
	log = zerolog.New(tempFile).With().Logger()
	log.Info().Msg("This is an initial entry from my log")
	fmt.Printf("The log file is allocated at %s\n", tempFile.Name())
	return &Logger{log: log}
}
