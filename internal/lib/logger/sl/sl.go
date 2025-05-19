// Package sl provides a flexible logging solution based on slog with support for multiple formats and outputs.
//
// Key features:
//   - Multiple output formats: JSON, pretty-printed, plain text, or discard
//   - File output with rotation (using lumberjack)
//   - Level-based filtering (DEBUG, INFO, WARN, ERROR)
//   - Functional options for configuration
//   - Error attribute helper
//
// Example usage:
//
//	logger := sl.MustLogger(
//		sl.WithLevel("debug"),
//		sl.WithFormat("pretty"),
//		sl.WithOutput("file"),
//		sl.WithFileOptions("./logs/app.log", 100, 7, true),
//	)
//	logger.Info("Application started", "version", "1.0.0")
package sl

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/10Narratives/distgo-db/internal/lib/logger/handlers/slogdiscard"
	"github.com/10Narratives/distgo-db/internal/lib/logger/handlers/slogpretty"
	"github.com/natefinch/lumberjack"
)

// Err creates a standardized error attribute for logging
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

// LoggerOptions contains configuration parameters for the logger
type LoggerOptions struct {
	level    slog.Level
	format   string
	output   string
	filePath string
	maxSize  int
	maxAge   int
	compress bool
}

// LoggerOption is a functional option for configuring the logger
type LoggerOption func(*LoggerOptions)

func defaultOptions() *LoggerOptions {
	return &LoggerOptions{
		level:    slog.LevelInfo,
		format:   "json",
		output:   "stdout",
		maxSize:  100,
		maxAge:   7,
		compress: true,
	}
}

// WithLevel sets the log level filter (debug/info/warn/error)
func WithLevel(levelStr string) LoggerOption {
	return func(o *LoggerOptions) {
		o.level = parseLogLevel(levelStr)
	}
}

// WithFormat sets the output format:
//   - "json": Structured JSON format
//   - "pretty": Colorized human-readable format
//   - "plain": Simple text format
//   - "discard": Disables logging
func WithFormat(format string) LoggerOption {
	return func(o *LoggerOptions) {
		o.format = format
	}
}

// WithOutput sets the log output destination ("stdout" or "file")
func WithOutput(output string) LoggerOption {
	return func(o *LoggerOptions) {
		o.output = output
	}
}

// WithFileOptions configures file output parameters
func WithFileOptions(filePath string, maxSize, maxAge int, compress bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.filePath = filePath
		o.maxSize = maxSize
		o.maxAge = maxAge
		o.compress = compress
	}
}

// MustLogger creates a new logger with given options (panics on error)
func MustLogger(opts ...LoggerOption) *slog.Logger {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	output := createOutput(options)
	handler := createHandler(options.format, output, options.level)
	return slog.New(handler)
}

func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		panic("unsupported log level: " + levelStr)
	}
}

// createOutput creates the appropriate io.Writer based on configuration
func createOutput(opts *LoggerOptions) io.Writer {
	if opts.output == "stdout" {
		return os.Stdout
	}

	dir := filepath.Dir(opts.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	return &lumberjack.Logger{
		Filename:  opts.filePath,
		MaxSize:   opts.maxSize,
		MaxAge:    opts.maxAge,
		Compress:  opts.compress,
		LocalTime: true,
	}
}

// createHandler creates the appropriate slog.Handler based on format
func createHandler(format string, output io.Writer, level slog.Level) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch format {
	case "json":
		return slog.NewJSONHandler(output, opts)
	case "pretty":
		return slogpretty.NewPrettyLogger(&slogpretty.PrettyHandlerOptions{
			SlogOpts: opts,
		}, output).Handler()
	case "plain":
		return slog.NewTextHandler(output, opts)
	case "discard":
		return slogdiscard.NewDiscardLogger().Handler()
	default:
		panic("unsupported log format: " + format)
	}
}
