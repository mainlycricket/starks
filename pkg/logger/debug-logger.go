package logger

import "log"

func (l *logger) DebugInfo(msg string, args ...any) {
	if l.debug {
		log.Println("info", msg, args)
	}
}

func (l *logger) DebugError(msg string, args ...any) {
	if l.debug {
		log.Println("error", msg, args)
	}
}
