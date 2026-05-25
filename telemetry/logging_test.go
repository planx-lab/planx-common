package telemetry

import (
	"context"
	"testing"
)

func TestInitLogging(t *testing.T) {
	err := InitLogging(context.Background(), LoggingConfig{
		ServiceName: "test-service",
	})
	if err != nil {
		t.Fatalf("InitLogging: %v", err)
	}
}

func TestGetLoggerProvider(t *testing.T) {
	_ = InitLogging(context.Background(), LoggingConfig{
		ServiceName: "test-service",
	})

	p := GetLoggerProvider()
	if p == nil {
		t.Fatal("GetLoggerProvider returned nil after InitLogging")
	}
}

func TestShutdownLogging(t *testing.T) {
	_ = InitLogging(context.Background(), LoggingConfig{
		ServiceName: "test-service",
	})

	err := ShutdownLogging(context.Background())
	if err != nil {
		t.Fatalf("ShutdownLogging: %v", err)
	}
}

func TestShutdownLogging_WithoutInit(t *testing.T) {
	// Should not panic when provider is nil
	err := ShutdownLogging(context.Background())
	if err != nil {
		t.Fatalf("ShutdownLogging without init: %v", err)
	}
}
