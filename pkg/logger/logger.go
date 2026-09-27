package logger

import (
	"os"
)

type logger struct {
	debug bool
}

var Logger = newLogger()

func newLogger() *logger {
	return &logger{
		debug: os.Getenv("DEBUG") == "TRUE",
	}
}
