package http

import (
	"go.opentelemetry.io/otel/metric"

	"github.com/prskr/go-dito/infrastructure/telemetry"
)

var (
	tracer                 = telemetry.Tracer("handlers/http")
	meter                  = telemetry.Meter("handlers/http")
	domainHandlerHistogram metric.Int64Histogram
)

func init() {
	var err error

	domainHandlerHistogram, err = meter.Int64Histogram(
		"http.domain_handler.duration",
		metric.WithDescription("Duration to process requests after the domain handler"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		panic(err)
	}
}
