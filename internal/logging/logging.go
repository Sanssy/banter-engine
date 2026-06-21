package logging

import (
	"fmt"
	"io"
	"log"
)

type Logger struct {
	component string
	logger    *log.Logger
}

func New(output io.Writer, component string) *Logger {
	return &Logger{
		component: component,
		logger:    log.New(output, "", log.Ldate|log.Ltime|log.LUTC),
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.log("INFO", format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.log("ERROR", format, args...)
}

func (l *Logger) log(level, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.logger.Printf("level=%s component=%s message=%q", level, l.component, message)
}
