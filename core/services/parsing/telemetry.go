package parsing

import (
	"github.com/prskr/go-dito/infrastructure/telemetry"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter           = telemetry.Meter("core/services/routing")
	oasRulesCounter metric.Int64Counter
)

func init() {
	var err error
	oasRulesCounter, err = meter.Int64Counter("routing.oas_rules")
	if err != nil {
		panic(err)
	}
}
