package core_module

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

const (
	logLevelConfigKey  = "log.level"
	logFormatConfigKey = "log.format"
)

var (
	LogLevel  slog.Level
	LogFormat string
)

func StringToSlogLevel(
	level string,
) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func setupLog() {
	LogLevel = StringToSlogLevel(Config.GetString(logLevelConfigKey))
	LogFormat = Config.GetString(logFormatConfigKey)

	var handler slog.Handler
	if LogFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: LogLevel})
		slog.SetDefault(slog.New(handler))
	} else {
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: LogLevel})
		slog.SetDefault(slog.New(handler))
	}
}
