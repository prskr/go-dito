package httpx

import (
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/prskr/go-dito/infrastructure/logging"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context()).SpanContext()

		logAttrs := []any{
			slog.String("method", r.Method),
			slog.String("host", r.Host),
			slog.String("path", r.RequestURI),
		}

		if span.HasSpanID() && span.HasTraceID() {
			logAttrs = append(
				logAttrs,
				slog.String("trace_id", span.TraceID().String()),
				slog.String("span_id", span.SpanID().String()),
			)
		}

		logger := logging.GetLogger(r.Context()).With(logAttrs...)

		logger.Info("Handling request")

		next.ServeHTTP(w, r.WithContext(logging.ContextWithLogger(r.Context(), logger)))

		logger.Debug("Request complete")
	})
}
