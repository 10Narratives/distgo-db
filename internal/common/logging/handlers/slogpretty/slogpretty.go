// Package slogpretty provides a custom slog.Handler that formats log records
// with colorized levels, indented JSON fields, and timestamp prefixes for
// human-readable output.
package slogpretty

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/fatih/color"
)

// PrettyHandlerOptions configures the behavior of the PrettyHandler.
// It embeds slog.HandlerOptions to inherit standard handler configurations.
type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

// PrettyHandler is a slog.Handler implementation that outputs logs in a
// human-friendly format with ANSI color codes and indented JSON fields.
type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

// NewPrettyHandler creates a new PrettyHandler that writes to the given output.
// It combines colorized level prefixes, timestamp, message, and indented JSON
// attributes for improved readability.
func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

// Handle processes a log record by formatting it with colorized level, timestamp,
// message, and pretty-printed JSON attributes. Implements slog.Handler.
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]any, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

// WithAttrs returns a new PrettyHandler with additional attributes prepended
// to future log records. Implements slog.Handler.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

// WithGroup returns a new PrettyHandler that starts a group for subsequent
// attributes. Implements slog.Handler.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
