package logger

import (
	"log/slog"
	"os"
)

// New creates a new logger
func New() *slog.Logger {
	// Configure the logger
	opts := &slog.HandlerOptions{
		Level: getLogLevel(),
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

// getLogLevel returns the log level based on the environment variable
func getLogLevel() slog.Level {
	level := os.Getenv("LOG_LEVEL")
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
