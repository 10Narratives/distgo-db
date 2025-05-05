// Package slogdiscard provides a no-op slog.Logger implementation that discards
// all log records. Useful for disabling logging in specific contexts.
package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger creates a slog.Logger that discards all log output.
// Equivalent to slog.New(slogdiscard.NewDiscardHandler()).
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler implements slog.Handler by discarding all log records.
// All methods are no-ops and return the receiver itself.
type DiscardHandler struct{}

// NewDiscardHandler creates a new DiscardHandler instance.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle implements slog.Handler by discarding the log record.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs implements slog.Handler by returning the same DiscardHandler.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup implements slog.Handler by returning the same DiscardHandler.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled implements slog.Handler by always returning false.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
