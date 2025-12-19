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

func TestSpanHelpers(t *testing.T) {
	ctx := context.Background()

	ctx1, span1 := StartSourceReadSpan(ctx, "tenant-1", "session-1", 100)
	span1.End()
	if ctx1 == nil {
		t.Error("StartSourceReadSpan returned nil context")
	}

	ctx2, span2 := StartProcessorSpan(ctx, "json", "session-1", 100)
	span2.End()
	if ctx2 == nil {
		t.Error("StartProcessorSpan returned nil context")
	}

	ctx3, span3 := StartSinkWriteSpan(ctx, "http", "session-1", 100)
	span3.End()
	if ctx3 == nil {
		t.Error("StartSinkWriteSpan returned nil context")
	}

	ctx4, span4 := StartRouteSpan(ctx, "source", "processor")
	span4.End()
	if ctx4 == nil {
		t.Error("StartRouteSpan returned nil context")
	}
}

func TestRecordSpanError(t *testing.T) {
	ctx := context.Background()
	_, span := StartSpan(ctx, "test-span")
	defer span.End()

	RecordSpanError(span, context.Canceled)
}
