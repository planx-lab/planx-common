package logger

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestInit(t *testing.T) {
	buf := &bytes.Buffer{}
	Init(Config{
		Level:       "debug",
		Pretty:      false,
		Output:      buf,
		ServiceName: "test-service",
	})

	Info().Msg("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("expected log to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "test-service") {
		t.Errorf("expected log to contain service name, got: %s", output)
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name  string
		logFn func()
	}{
		{"Debug", func() { Debug().Msg("debug") }},
		{"Info", func() { Info().Msg("info") }},
		{"Warn", func() { Warn().Msg("warn") }},
		{"Error", func() { Error().Msg("error") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFn() // Should not panic
		})
	}
}

func TestWithContext(t *testing.T) {
	// Setup OTel tracer
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)
	tracer := provider.Tracer("test")

	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	l := WithContext(ctx)
	if l == nil {
		t.Fatal("WithContext returned nil")
	}
}

func TestContextAwareLogging(t *testing.T) {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)
	tracer := provider.Tracer("test")

	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	// These should not panic
	DebugCtx(ctx).Msg("debug with context")
	InfoCtx(ctx).Msg("info with context")
	WarnCtx(ctx).Msg("warn with context")
	ErrorCtx(ctx).Msg("error with context")
}

func TestAddSpanEvent(t *testing.T) {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)
	tracer := provider.Tracer("test")

	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	// Should not panic
	AddSpanEvent(ctx, "test event")
}
