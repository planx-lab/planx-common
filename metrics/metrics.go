// Package metrics provides metrics interfaces for Planx.
// This is an abstraction layer to allow different metrics backends.
package metrics

import "context"

// Counter represents a monotonically increasing counter.
type Counter interface {
	Inc()
	Add(delta float64)
}

// Gauge represents a value that can go up and down.
type Gauge interface {
	Set(value float64)
	Inc()
	Dec()
	Add(delta float64)
	Sub(delta float64)
}

// Histogram represents a distribution of values.
type Histogram interface {
	Observe(value float64)
}

// Provider is the interface for metrics providers.
type Provider interface {
	// Counter returns a counter with the given name and labels.
	Counter(name string, labels map[string]string) Counter

	// Gauge returns a gauge with the given name and labels.
	Gauge(name string, labels map[string]string) Gauge

	// Histogram returns a histogram with the given name and labels.
	Histogram(name string, labels map[string]string) Histogram
}

// Recorder provides high-level metrics recording.
type Recorder interface {
	// RecordBatchProcessed records a batch was processed.
	RecordBatchProcessed(ctx context.Context, pluginName string, recordCount int)

	// RecordBatchLatency records the latency of processing a batch.
	RecordBatchLatency(ctx context.Context, pluginName string, latencyMs float64)

	// RecordSessionActive records the number of active sessions.
	RecordSessionActive(ctx context.Context, pluginName string, count int)

	// RecordError records an error occurrence.
	RecordError(ctx context.Context, pluginName string, errorType string)
}

// NoopCounter is a no-op counter for testing.
type NoopCounter struct{}

func (NoopCounter) Inc()          {}
func (NoopCounter) Add(_ float64) {}

// NoopGauge is a no-op gauge for testing.
type NoopGauge struct{}

func (NoopGauge) Set(_ float64) {}
func (NoopGauge) Inc()          {}
func (NoopGauge) Dec()          {}
func (NoopGauge) Add(_ float64) {}
func (NoopGauge) Sub(_ float64) {}

// NoopHistogram is a no-op histogram for testing.
type NoopHistogram struct{}

func (NoopHistogram) Observe(_ float64) {}

// NoopProvider is a no-op metrics provider.
type NoopProvider struct{}

func (NoopProvider) Counter(_ string, _ map[string]string) Counter     { return NoopCounter{} }
func (NoopProvider) Gauge(_ string, _ map[string]string) Gauge         { return NoopGauge{} }
func (NoopProvider) Histogram(_ string, _ map[string]string) Histogram { return NoopHistogram{} }
