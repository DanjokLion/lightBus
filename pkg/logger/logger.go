package logger

import (
	"log"
)

type Logger struct {
	prefix string
}

func New() *Logger {
	return &Logger{
		prefix: "[lightbus]",
	}
}

// todo: replace Printf to zap/zerolog 

func (l *Logger) Info(format string, v ...interface{}) {
	log.Printf(l.prefix + " INFO: "+ format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	log.Printf(l.prefix + " ERROR: " + format, v...)
}