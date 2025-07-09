package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// NewConsoleHandler creates a custom console handler (simplified)
func NewConsoleHandler(w io.Writer, opts *slog.HandlerOptions, enableColors bool) slog.Handler {
	return slog.NewTextHandler(w, opts)
}

// Default returns a default logger instance
func Default() *Logger {
	return New(DefaultConfig())
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// loggerKey is the key for storing/retrieving logger in context
const loggerKey contextKey = "logger"

// FromContext extracts logger from context
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	return Default()
}

// WithLogger adds logger to context
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
