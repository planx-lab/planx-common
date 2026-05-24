package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
	tracingOnce    sync.Once
	tracerInit     sync.Once
)

// TracingConfig holds tracing configuration.
type TracingConfig struct {
	ServiceName string
	Endpoint    string // OTLP endpoint, empty for stdout
}

// InitTracing initializes OpenTelemetry tracing.
func InitTracing(ctx context.Context, cfg TracingConfig) error {
	var initErr error
	tracingOnce.Do(func() {
		initErr = initTracingInternal(ctx, cfg)
	})
	return initErr
}

func initTracingInternal(ctx context.Context, cfg TracingConfig) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return err
	}

	var exporter sdktrace.SpanExporter
	if cfg.Endpoint != "" {
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(cfg.Endpoint),
		}
		exporter, err = otlptrace.New(ctx, otlptracehttp.NewClient(opts...))
	} else {
		exporter, err = stdouttrace.New()
	}
	if err != nil {
		return err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracerProvider = provider
	tracer = provider.Tracer("planx")

	return nil
}

// ShutdownTracing gracefully shuts down the tracer provider.
func ShutdownTracing(ctx context.Context) error {
	if tracerProvider != nil {
		tp := tracerProvider
		tracerProvider = nil
		return tp.Shutdown(ctx)
	}
	return nil
}

// Tracer returns the global tracer.
func Tracer() trace.Tracer {
	tracerInit.Do(func() {
		if tracer == nil {
			tracer = otel.Tracer("planx")
		}
	})
	return tracer
}

// StartSpan starts a new span with the given name.
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, name, trace.WithAttributes(attrs...))
}

// SpanFromContext returns the current span from context.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceID returns the trace ID from context.
func TraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// SpanID returns the span ID from context.
func SpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasSpanID() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// InjectTraceContext injects trace context into a map (for Batch.Context).
func InjectTraceContext(ctx context.Context, carrier map[string]string) {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.MapCarrier(carrier))
}

// ExtractTraceContext extracts trace context from a map.
func ExtractTraceContext(ctx context.Context, carrier map[string]string) context.Context {
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, propagation.MapCarrier(carrier))
}

// Span helpers for common operations

// StartSourceReadSpan starts a span for source.read.
func StartSourceReadSpan(ctx context.Context, tenantID, sessionID string, batchSize int) (context.Context, trace.Span) {
	return StartSpan(ctx, "source.read",
		attribute.String("tenant_id", tenantID),
		attribute.String("session_id", sessionID),
		attribute.Int("batch_size", batchSize),
	)
}

// StartProcessorSpan starts a span for processor.process.
func StartProcessorSpan(ctx context.Context, processorName, sessionID string, batchSize int) (context.Context, trace.Span) {
	return StartSpan(ctx, "processor.process",
		attribute.String("processor", processorName),
		attribute.String("session_id", sessionID),
		attribute.Int("batch_size", batchSize),
	)
}

// StartSinkWriteSpan starts a span for sink.write.
func StartSinkWriteSpan(ctx context.Context, sinkName, sessionID string, batchSize int) (context.Context, trace.Span) {
	return StartSpan(ctx, "sink.write",
		attribute.String("sink", sinkName),
		attribute.String("session_id", sessionID),
		attribute.Int("batch_size", batchSize),
	)
}

// StartRouteSpan starts a span for engine.route.
func StartRouteSpan(ctx context.Context, fromStage, toStage string) (context.Context, trace.Span) {
	return StartSpan(ctx, "engine.route",
		attribute.String("from", fromStage),
		attribute.String("to", toStage),
	)
}

// RecordSpanError records an error on the current span.
func RecordSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
