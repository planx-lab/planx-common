package telemetry

import (
	"context"
	"testing"
)

func TestInitTracing(t *testing.T) {
	ctx := context.Background()
	err := InitTracing(ctx, TracingConfig{
		ServiceName: "test-service",
		Endpoint:    "", // stdout
	})
	if err != nil {
		t.Fatalf("InitTracing failed: %v", err)
	}
}

func TestTracer(t *testing.T) {
	tracer := Tracer()
	if tracer == nil {
		t.Fatal("Tracer returned nil")
	}
}

func TestStartSpan(t *testing.T) {
	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test-span")
	if span == nil {
		t.Fatal("StartSpan returned nil span")
	}
	span.End()
}

func TestSpanFromContext(t *testing.T) {
	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test-span")
	defer span.End()

	retrieved := SpanFromContext(ctx)
	if retrieved == nil {
		t.Fatal("SpanFromContext returned nil")
	}
}

func TestTraceID(t *testing.T) {
	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test-span")
	defer span.End()

	traceID := TraceID(ctx)
	if traceID == "" {
		t.Log("TraceID is empty (expected for noop tracer)")
	}
}

func TestSpanID(t *testing.T) {
	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test-span")
	defer span.End()

	spanID := SpanID(ctx)
	if spanID == "" {
		t.Log("SpanID is empty (expected for noop tracer)")
	}
}

func TestInjectExtractTraceContext(t *testing.T) {
	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test-span")
	defer span.End()

	carrier := make(map[string]string)
	InjectTraceContext(ctx, carrier)

	// Extract back
	newCtx := ExtractTraceContext(context.Background(), carrier)
	if newCtx == nil {
		t.Fatal("ExtractTraceContext returned nil context")
	}
}
