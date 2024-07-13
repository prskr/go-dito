package telemetry

import (
	"path"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func Meter(name string, opts ...metric.MeterOption) metric.Meter {
	initBuildInfo()
	if infoIsSet {
		return otel.Meter(path.Join(buildInfo.Main.Path, name), opts...)
	}
	return otel.Meter(name, opts...)
}
