package logger

import (
	"log/slog"
	"os"
)

// New membuat structured logger menggunakan stdlib slog.
// Development: text format yang mudah dibaca di terminal.
// Production: JSON format untuk log aggregation (ELK, Loki, dll).
func New(env string) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if env == "development" {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler).With(
		slog.String("service", "healy-backend"),
		slog.String("env", env),
	)
}
