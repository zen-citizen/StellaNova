package utils

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type TraceHandler struct {
	slog.Handler
}

func NewTraceHandler(h slog.Handler) *TraceHandler {
	return &TraceHandler{
		Handler: h,
	}
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()))
	}
	return h.Handler.Handle(ctx, r)
}

func SetupLogger(level, format string) *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	var slogHandler slog.Handler
	if strings.ToLower(format) == "json" {
		slogHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		slogHandler = slog.NewTextHandler(os.Stdout, opts)
	}

	traceHandler := NewTraceHandler(slogHandler)
	return slog.New(traceHandler)
}
