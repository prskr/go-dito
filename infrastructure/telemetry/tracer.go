package telemetry

import (
	"path"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	initBuildInfo()
	if infoIsSet {
		return otel.GetTracerProvider().Tracer(path.Join(buildInfo.Main.Path, name), opts...)
	}
	return otel.GetTracerProvider().Tracer(name, opts...)
}
