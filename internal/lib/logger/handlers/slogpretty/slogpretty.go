// Package slogpretty provides a human-readable, colorized log formatter for the slog package.
package slogpretty

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

// PrettyHandlerOptions contains configuration options for the pretty printer
type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

// PrettyHandler implements a human-readable slog.Handler with colorized output
type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

// NewPrettyHandler creates a new PrettyHandler instance
// Parameters:
//   - out: io.Writer to send formatted log output to (typically os.Stdout)
//
// Returns a handler that:
//   - Colorizes log levels
//   - Formats timestamps as [HH:MM:SS.MSS]
//   - Pretty-prints attributes as indented JSON
//   - Maintains attribute context with WithAttrs
func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

// Handle formats and writes a log record
// Format: [TIME] LEVEL: MESSAGE {ATTRS}
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

// WithAttrs adds contextual attributes to the handler
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

// WithGroup starts a new attribute group
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}

// NewPrettyLogger creates a configured slog.Logger
// Parameters:
//   - opts: Configuration options (use nil for defaults)
//   - out: Optional output destination (defaults to os.Stdout)
//
// Returns a logger that:
//   - Uses colorized output for levels
//   - Shows millisecond timestamps
//   - Pretty-prints attributes
//   - Filters logs by level (default: Info)
func NewPrettyLogger(opts *PrettyHandlerOptions, out ...io.Writer) *slog.Logger {
	output := io.Writer(os.Stdout)
	if len(out) > 0 {
		output = out[0]
	}

	if opts == nil {
		opts = &PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		}
	}

	handler := opts.NewPrettyHandler(output)
	return slog.New(handler)
}
