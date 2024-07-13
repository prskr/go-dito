package services

import (
	"go.opentelemetry.io/otel/metric"

	"github.com/prskr/go-dito/infrastructure/telemetry"
)

var (
	meter           = telemetry.Meter("infrastructure/telemetry")
	oasRulesCounter metric.Int64Counter
)

func init() {
	var err error
	oasRulesCounter, err = meter.Int64Counter("routing.oas_rules")
	if err != nil {
		panic(err)
	}
}
