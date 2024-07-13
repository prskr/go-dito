package config

import (
	"io"
	"log/slog"
	"strings"
)

func (l Logging) Handler(out io.Writer) slog.Handler {
	switch l.Format.String() {
	case "text":
		return slog.NewTextHandler(out, l.Options())
	case "json":
		fallthrough
	default:
		return slog.NewJSONHandler(out, l.Options())
	}
}

func (l Logging) Options() *slog.HandlerOptions {
	opts := &slog.HandlerOptions{
		AddSource: l.AddSource,
	}

	switch strings.ToLower(l.Level.String()) {
	case "debug":
		opts.Level = slog.LevelDebug
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	case "info":
		fallthrough
	default:
		opts.Level = slog.LevelInfo
	}

	return opts
}
