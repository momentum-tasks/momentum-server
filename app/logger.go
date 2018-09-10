package app

import (
	"os"

	"github.com/op/go-logging"
)

// Logger type is not fully utilized, but will be expanded to include the ability to have multiple instances to log to different locations
type Logger struct {
}

var mlog = logging.MustGetLogger("momentum")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
)

// NewLogger creates the logger object to be referenced throughout the application
func NewLogger() *Logger {
	return new(Logger)
}

// Begin creates and backends of the logger, and prepares them for use
func (l *Logger) Begin() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(logBackend, format)
	logging.SetBackend(backendFormatter)
}
