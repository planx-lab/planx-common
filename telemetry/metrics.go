// Package telemetry provides OpenTelemetry instrumentation for Planx.
// Engine-side utilities only — must not be imported by SDK or plugins.
package telemetry

import (
	"context"
	"errors"
	"fmt"
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
	if err := initInstruments(provider); err != nil {
		return err
	}

	return nil
}

// InitMetricsWithReaders initializes metrics with custom readers.
// This allows callers (like the engine) to provide their own reader
// (e.g., Prometheus ManualReader) instead of using the default
// PeriodicReader + OTLP/stdout exporter.
// Returns the MeterProvider for lifecycle management.
// Must be called only once in production; tests may call it per-test.
func InitMetricsWithReaders(ctx context.Context, cfg MetricsConfig, readers ...sdkmetric.Reader) (*sdkmetric.MeterProvider, error) {
	return initMetricsWithReadersInternal(ctx, cfg, readers...)
}

func initMetricsWithReadersInternal(ctx context.Context, cfg MetricsConfig, readers ...sdkmetric.Reader) (*sdkmetric.MeterProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	opts := []sdkmetric.Option{sdkmetric.WithResource(res)}
	for _, r := range readers {
		opts = append(opts, sdkmetric.WithReader(r))
	}

	provider := sdkmetric.NewMeterProvider(opts...)
	otel.SetMeterProvider(provider)
	if err := initInstruments(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func initInstruments(provider *sdkmetric.MeterProvider) error {
	meter = provider.Meter("planx")

	var errs []error

	// Initialize instruments
	var err error
	batchesSent, err = meter.Int64Counter("planx.batches.sent",
		metric.WithDescription("Total batches sent"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating batches.sent counter: %w", err))
	}
	batchesReceived, err = meter.Int64Counter("planx.batches.received",
		metric.WithDescription("Total batches received"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating batches.received counter: %w", err))
	}
	recordsSent, err = meter.Int64Counter("planx.records.sent",
		metric.WithDescription("Total records sent"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating records.sent counter: %w", err))
	}
	recordsReceived, err = meter.Int64Counter("planx.records.received",
		metric.WithDescription("Total records received"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating records.received counter: %w", err))
	}
	errorsTotal, err = meter.Int64Counter("planx.errors.total",
		metric.WithDescription("Total errors"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating errors.total counter: %w", err))
	}

	stageLatency, err = meter.Float64Histogram("planx.stage.latency",
		metric.WithDescription("Stage processing latency in milliseconds"),
		metric.WithUnit("ms"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating stage.latency histogram: %w", err))
	}
	ackLatency, err = meter.Float64Histogram("planx.ack.latency",
		metric.WithDescription("ACK latency in milliseconds"),
		metric.WithUnit("ms"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating ack.latency histogram: %w", err))
	}

	windowBacklog, err = meter.Int64UpDownCounter("planx.window.backlog",
		metric.WithDescription("Window backlog (in-flight batches)"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating window.backlog updowncounter: %w", err))
	}
	sessionsActive, err = meter.Int64UpDownCounter("planx.sessions.active",
		metric.WithDescription("Active sessions"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating sessions.active updowncounter: %w", err))
	}
	inFlightBatches, err = meter.Int64UpDownCounter("planx.batches.inflight",
		metric.WithDescription("In-flight batches"))
	if err != nil {
		errs = append(errs, fmt.Errorf("creating batches.inflight updowncounter: %w", err))
	}

	return errors.Join(errs...)
}

// RecordBatchSent records a batch being sent.
func RecordBatchSent(ctx context.Context, tenantID, stage, pluginType string, recordCount int64) {
	if batchesSent == nil || recordsSent == nil {
		return
	}
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
	if batchesReceived == nil || recordsReceived == nil {
		return
	}
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
	if stageLatency == nil {
		return
	}
	stageLatency.Record(ctx, latencyMs, metric.WithAttributes(
		attribute.String("stage", stage),
	))
}

// RecordAckLatency records the ACK latency.
func RecordAckLatency(ctx context.Context, latencyMs float64) {
	if ackLatency == nil {
		return
	}
	ackLatency.Record(ctx, latencyMs)
}

// RecordError records an error.
func RecordError(ctx context.Context, tenantID, stage, errorType string) {
	if errorsTotal == nil {
		return
	}
	errorsTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.String("stage", stage),
		attribute.String("error_type", errorType),
	))
}

// UpdateWindowBacklog updates the window backlog gauge.
func UpdateWindowBacklog(ctx context.Context, stage string, delta int64) {
	if windowBacklog == nil {
		return
	}
	windowBacklog.Add(ctx, delta, metric.WithAttributes(
		attribute.String("stage", stage),
	))
}

// UpdateSessionsActive updates the active sessions gauge.
func UpdateSessionsActive(ctx context.Context, pluginType string, delta int64) {
	if sessionsActive == nil {
		return
	}
	sessionsActive.Add(ctx, delta, metric.WithAttributes(
		attribute.String("plugin_type", pluginType),
	))
}

// UpdateInFlightBatches updates the in-flight batches gauge.
func UpdateInFlightBatches(ctx context.Context, delta int64) {
	if inFlightBatches == nil {
		return
	}
	inFlightBatches.Add(ctx, delta)
}
