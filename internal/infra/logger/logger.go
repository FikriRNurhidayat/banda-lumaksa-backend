package logger

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/fikrirnurhidayat/banda-lumaksa/pkg/exists"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
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
var Any = slog.Any

var levels = map[string]slog.Level{
	"error": slog.LevelError,
	"warn":  slog.LevelWarn,
	"info":  slog.LevelInfo,
	"debug": slog.LevelDebug,
}

func New() Logger {
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.source", false)
	viper.SetDefault("log.time", true)

	bi, _ := debug.ReadBuildInfo()

	level, ok := levels[viper.GetString("log.level")]
	if !ok {
		level = levels["info"]
	}

	var handler slog.Handler
	isTimeLogged := viper.GetBool("log.time")
	opts := &slog.HandlerOptions{
		AddSource: viper.GetBool("log.source"),
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && !isTimeLogged {
				return slog.Attr{}
			}

			return a
		},
	}

	switch viper.GetString("log.style") {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	if exists.String(bi.Main.Version) {
		logger = logger.With(slog.String("v", bi.Main.Version))
	}

	return logger
}
