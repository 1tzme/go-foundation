package logger

import (
	"log/slog"
	"os"
)

// Config holds logger configuration
type Config struct {
	Level  string
	Format string // "json" or "text"
}

// Logger wraps slog.Logger with additional context
type Logger struct {
	*slog.Logger
}

// New creates a new logger instance with the given configuration
func New(config Config) *Logger {
	var level slog.Level
	switch config.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	return &Logger{Logger: logger}
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(args ...interface{}) *Logger {
	return &Logger{Logger: l.Logger.With(args...)}
}

// Default returns the default logger instance
func Default() *Logger {
	return &Logger{Logger: slog.Default()}
}
