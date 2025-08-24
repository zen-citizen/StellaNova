package utils

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log/slog"
)

func SetupTracing(serviceName string) func() {
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		slog.Error("failed to initialize resource", slog.Any("error", err))
		return func() {}
	}

	tp := trace.NewTracerProvider(
		trace.WithResource(r),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown exporter", slog.Any("error", err))
		}
	}
}
