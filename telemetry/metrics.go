// Package otel provides OpenTelemetry instrumentation for Planx.
// This package is shared by engine and plugins.
package telemetry

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	meter     metric.Meter
	meterOnce sync.Once

	// Counters
	batchesSent     metric.Int64Counter
	batchesReceived metric.Int64Counter
	recordsSent     metric.Int64Counter
	recordsReceived metric.Int64Counter
	errorsTotal     metric.Int64Counter

	// Histograms
	stageLatency metric.Float64Histogram
	ackLatency   metric.Float64Histogram

	// Gauges
	windowBacklog   metric.Int64UpDownCounter
	sessionsActive  metric.Int64UpDownCounter
	inFlightBatches metric.Int64UpDownCounter
)

// MetricsConfig holds metrics configuration.
type MetricsConfig struct {
	ServiceName string
	Endpoint    string // OTLP endpoint, empty for stdout
	Interval    time.Duration
}

// InitMetrics initializes OpenTelemetry metrics.
func InitMetrics(ctx context.Context, cfg MetricsConfig) error {
	var err error
	meterOnce.Do(func() {
		err = initMetricsInternal(ctx, cfg)
	})
	return err
}

func initMetricsInternal(ctx context.Context, cfg MetricsConfig) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return err
	}

	var exporter sdkmetric.Exporter
	if cfg.Endpoint != "" {
		exporter, err = otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
			otlpmetricgrpc.WithInsecure(),
		)
	} else {
		exporter, err = stdoutmetric.New()
	}
	if err != nil {
		return err
	}

	interval := cfg.Interval
	if interval == 0 {
		interval = 10 * time.Second
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(interval))),
	)

	otel.SetMeterProvider(provider)
	meter = provider.Meter("planx")

	// Initialize instruments
	batchesSent, _ = meter.Int64Counter("planx.batches.sent",
		metric.WithDescription("Total batches sent"))
	batchesReceived, _ = meter.Int64Counter("planx.batches.received",
		metric.WithDescription("Total batches received"))
	recordsSent, _ = meter.Int64Counter("planx.records.sent",
		metric.WithDescription("Total records sent"))
	recordsReceived, _ = meter.Int64Counter("planx.records.received",
		metric.WithDescription("Total records received"))
	errorsTotal, _ = meter.Int64Counter("planx.errors.total",
		metric.WithDescription("Total errors"))

	stageLatency, _ = meter.Float64Histogram("planx.stage.latency",
		metric.WithDescription("Stage processing latency in milliseconds"),
		metric.WithUnit("ms"))
	ackLatency, _ = meter.Float64Histogram("planx.ack.latency",
		metric.WithDescription("ACK latency in milliseconds"),
		metric.WithUnit("ms"))

	windowBacklog, _ = meter.Int64UpDownCounter("planx.window.backlog",
		metric.WithDescription("Window backlog (in-flight batches)"))
	sessionsActive, _ = meter.Int64UpDownCounter("planx.sessions.active",
		metric.WithDescription("Active sessions"))
	inFlightBatches, _ = meter.Int64UpDownCounter("planx.batches.inflight",
		metric.WithDescription("In-flight batches"))

	return nil
}

// RecordBatchSent records a batch being sent.
func RecordBatchSent(ctx context.Context, tenantID, stage, pluginType string, recordCount int64) {
	attrs := []attribute.KeyValue{
		attribute.String("tenant_id", tenantID),
		attribute.String("stage", stage),
		attribute.String("plugin_type", pluginType),
	}
	batchesSent.Add(ctx, 1, metric.WithAttributes(attrs...))
	recordsSent.Add(ctx, recordCount, metric.WithAttributes(attrs...))
}

// RecordBatchReceived records a batch being received.
func RecordBatchReceived(ctx context.Context, tenantID, stage, pluginType string, recordCount int64) {
	attrs := []attribute.KeyValue{
		attribute.String("tenant_id", tenantID),
		attribute.String("stage", stage),
		attribute.String("plugin_type", pluginType),
	}
	batchesReceived.Add(ctx, 1, metric.WithAttributes(attrs...))
	recordsReceived.Add(ctx, recordCount, metric.WithAttributes(attrs...))
}

// RecordStageLatency records the latency for a pipeline stage.
func RecordStageLatency(ctx context.Context, stage string, latencyMs float64) {
	stageLatency.Record(ctx, latencyMs, metric.WithAttributes(
		attribute.String("stage", stage),
	))
}

// RecordAckLatency records the ACK latency.
func RecordAckLatency(ctx context.Context, latencyMs float64) {
	ackLatency.Record(ctx, latencyMs)
}

// RecordError records an error.
func RecordError(ctx context.Context, tenantID, stage, errorType string) {
	errorsTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.String("stage", stage),
		attribute.String("error_type", errorType),
	))
}

// UpdateWindowBacklog updates the window backlog gauge.
func UpdateWindowBacklog(ctx context.Context, stage string, delta int64) {
	windowBacklog.Add(ctx, delta, metric.WithAttributes(
		attribute.String("stage", stage),
	))
}

// UpdateSessionsActive updates the active sessions gauge.
func UpdateSessionsActive(ctx context.Context, pluginType string, delta int64) {
	sessionsActive.Add(ctx, delta, metric.WithAttributes(
		attribute.String("plugin_type", pluginType),
	))
}

// UpdateInFlightBatches updates the in-flight batches gauge.
func UpdateInFlightBatches(ctx context.Context, delta int64) {
	inFlightBatches.Add(ctx, delta)
}
