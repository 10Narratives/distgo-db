// Package slogdiscard provides a no-op implementation of slog.Handler.
// It's useful for disabling logging in tests or specific environments.
package slogdiscard

import (
	"context"
	"log/slog"
)

// DiscardHandler is a slog.Handler implementation that discards all log records.
type DiscardHandler struct{}

// NewDiscardHandler creates a new DiscardHandler instance.
// This handler will ignore all log messages regardless of their level.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Enabled always returns false to disable all log levels.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

// Handle discards log records by doing nothing.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the same handler since no attributes need to be stored.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the same handler since no groups need to be tracked.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// NewDiscardLogger creates a slog.Logger that discards all output.
// Useful for testing or when you need to completely disable logging.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}
