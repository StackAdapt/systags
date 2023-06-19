package command

import (
	"golang.org/x/exp/slog"
)

var logger = slog.Default()

func GetLogger() *slog.Logger {
	return logger
}

func SetLogger(l *slog.Logger) {
	logger = l
}
