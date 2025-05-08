package logging

import (
	"context"
	"log/slog"
)

var loggerKey = struct {
	key string
}{
	key: "logger",
}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	contextLogger := ctx.Value(loggerKey).(*slog.Logger)
	if contextLogger == nil {
		return slog.Default()
	}
	return contextLogger
}
