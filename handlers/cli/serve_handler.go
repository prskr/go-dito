package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/apple/pkl-go/pkl"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/prskr/go-dito/core/services"
	http2 "github.com/prskr/go-dito/handlers/http"
	"github.com/prskr/go-dito/infrastructure/config"
	"github.com/prskr/go-dito/infrastructure/httpx"
	"github.com/prskr/go-dito/infrastructure/logging"
)

type ServeHandler struct {
}

func (h *ServeHandler) Run(ctx context.Context, cfg *config.AppConfig, logger *slog.Logger) error {
	domainHandler := make(http2.DomainHandler)

	for d, a := range cfg.Domains {
		parser, err := services.DefaultRegistry.ParserFor(a)
		if err != nil {
			return err
		}

		if handler, err := parser.Handler(ctx); err != nil {
			return err
		} else {
			domainHandler[d] = handler
		}
	}

	srv := http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		ReadHeaderTimeout: cfg.Server.ServerOptions.ReadHeaderTimeout.GoDuration(),
		Handler:           otelhttp.NewHandler(httpx.LoggingMiddleware(http.MaxBytesHandler(domainHandler, int64(cfg.Server.RequestOptions.MaxBodySize.ToUnit(pkl.Bytes).Value))), "API"),
		BaseContext: func(listener net.Listener) context.Context {
			return logging.ContextWithLogger(ctx, logger)
		},
	}

	slog.Info("Starting server", slog.String("addr", srv.Addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to listen and serve", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	shutdownCtx, stop := context.WithTimeout(context.Background(), cfg.Server.ServerOptions.ShutdownTimeout.GoDuration())
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	stop()

	return nil
}

func (h *ServeHandler) AfterApply(ctx context.Context, appCfg *config.AppConfig) error {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("go-dito"),
	)

	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create metric reader: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricReader),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	go func() {
		<-ctx.Done()
		shutdownCtx, stop := context.WithTimeout(context.Background(), appCfg.Telemetry.ShutdownTimeout.GoDuration())
		err := meterProvider.Shutdown(shutdownCtx)
		stop()

		if err != nil {
			slog.Error("Failed to shutdown metric provider", logging.Error(err))
		}
	}()

	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create span exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(spanExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(traceProvider)

	go func() {
		<-ctx.Done()
		shutdownCtx, stop := context.WithTimeout(context.Background(), appCfg.Telemetry.ShutdownTimeout.GoDuration())
		err := traceProvider.Shutdown(shutdownCtx)
		stop()

		if err != nil {
			slog.Error("Failed to shutdown trace provider", logging.Error(err))
		}
	}()

	return nil
}
