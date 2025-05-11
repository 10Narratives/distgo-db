package logging

import (
	"log/slog"
	"os"

	"githib.com/10Narratives/distgo-db/pkg/logging/handlers/slogpretty"
)

// TODO: refactor this file

func fetchSlogLevel(level string) slog.Level {
	var ll slog.Level
	switch level {
	case "info":
		ll = slog.LevelInfo
		break
	case "warn":
		ll = slog.LevelWarn
		break
	case "debug":
		ll = slog.LevelDebug
		break
	case "error":
		ll = slog.LevelError
		break
	}
	return ll
}

func NewLogger(format, level string) *slog.Logger {
	var log *slog.Logger

	ll := fetchSlogLevel(level)

	switch format {
	case "json":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: ll}))
		break
	case "text":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: ll}))
		break
	case "pretty":
		log = setupPrettyLog(ll)
		break
	}

	return log
}

func setupPrettyLog(ll slog.Level) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: ll,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
