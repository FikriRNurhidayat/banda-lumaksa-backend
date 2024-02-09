package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

var String = slog.String
var Int = slog.Int
var Int64 = slog.Int64

func New() Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
