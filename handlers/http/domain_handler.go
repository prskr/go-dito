package http

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var _ http.Handler = (*DomainHandler)(nil)

type DomainHandler map[string]http.Handler

func (d DomainHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now().UTC()
	ctx, span := tracer.Start(
		request.Context(),
		"SelectDomainHandler",
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(attribute.String("domain", request.Host)),
	)
	defer func() {
		span.End()
		domainHandlerHistogram.Record(
			ctx,
			time.Since(start).Milliseconds(),
			metric.WithAttributes(attribute.String("domain", request.Host)),
		)
	}()

	request = request.WithContext(ctx)

	h, ok := d[request.Host]
	if !ok {
		span.AddEvent("DomainNotFound")
		http.NotFound(writer, request)
		return
	}

	h.ServeHTTP(writer, request)
}
