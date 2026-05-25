package logger

import (
	"bytes"
	"context"
	"os"
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

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Level != "info" {
		t.Fatalf("Level: got %q, want %q", cfg.Level, "info")
	}
	if cfg.Pretty {
		t.Fatal("Pretty should be false by default")
	}
	if cfg.ServiceName != "planx" {
		t.Fatalf("ServiceName: got %q", cfg.ServiceName)
	}
	if cfg.Output != os.Stdout {
		t.Fatal("Output should be os.Stdout")
	}
}

func TestContextWithTrace(t *testing.T) {
	ctx := context.Background()
	ctx = ContextWithTrace(ctx, "trace-123", "span-456")

	if ctx.Value(traceIDKey).(string) != "trace-123" {
		t.Fatal("trace_id not set in context")
	}
	if ctx.Value(spanIDKey).(string) != "span-456" {
		t.Fatal("span_id not set in context")
	}
}

func TestAddSpanEventWithAttrs(t *testing.T) {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)
	tracer := provider.Tracer("test")

	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	AddSpanEventWithAttrs(ctx, "event msg", map[string]string{
		"key1": "val1",
		"key2": "val2",
	})
}

func TestAddSpanEventWithAttrs_NoSpan(t *testing.T) {
	// Should not panic when no active span
	AddSpanEventWithAttrs(context.Background(), "no span", map[string]string{"k": "v"})
}

func TestAddSpanEvent_NoSpan(t *testing.T) {
	// Should not panic when no active span
	AddSpanEvent(context.Background(), "no span")
}

func TestGet_AutoInit(t *testing.T) {
	// Get() auto-initializes with defaults if Init not called.
	// Since sync.Once is already triggered in this test binary,
	// we verify it doesn't panic.
	l := Get()
	if l == nil {
		t.Fatal("Get returned nil")
	}
}
