package routing

import "github.com/prskr/go-dito/infrastructure/telemetry"

var (
	tracer = telemetry.Tracer("infrastructure/routing")
	meter  = telemetry.Meter("infrastructure/routing")
)
