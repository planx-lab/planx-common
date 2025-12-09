package otel

import (
	"context"
	"testing"
)

func TestInitMetrics(t *testing.T) {
	ctx := context.Background()
	err := InitMetrics(ctx, MetricsConfig{
		ServiceName: "test-service",
		Endpoint:    "", // stdout
	})
	// May error due to once.Do, but should not panic
	_ = err
}

func TestRecordBatchSent(t *testing.T) {
	ctx := context.Background()
	// Should not panic even if metrics not initialized
	RecordBatchSent(ctx, "tenant-1", "source", "mysql", 100)
}

func TestRecordBatchReceived(t *testing.T) {
	ctx := context.Background()
	RecordBatchReceived(ctx, "tenant-1", "sink", "http", 100)
}

func TestRecordStageLatency(t *testing.T) {
	ctx := context.Background()
	RecordStageLatency(ctx, "processor", 5.5)
}

func TestRecordAckLatency(t *testing.T) {
	ctx := context.Background()
	RecordAckLatency(ctx, 2.5)
}

func TestRecordError(t *testing.T) {
	ctx := context.Background()
	RecordError(ctx, "tenant-1", "sink", "connection_refused")
}

func TestUpdateWindowBacklog(t *testing.T) {
	ctx := context.Background()
	UpdateWindowBacklog(ctx, "processor-1", 5)
	UpdateWindowBacklog(ctx, "processor-1", -2)
}

func TestUpdateSessionsActive(t *testing.T) {
	ctx := context.Background()
	UpdateSessionsActive(ctx, "source", 1)
	UpdateSessionsActive(ctx, "source", -1)
}

func TestUpdateInFlightBatches(t *testing.T) {
	ctx := context.Background()
	UpdateInFlightBatches(ctx, 10)
	UpdateInFlightBatches(ctx, -5)
}
