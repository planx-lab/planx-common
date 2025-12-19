// Package otel provides OpenTelemetry logging integration for Planx.
// This enables log export via OTLP alongside metrics and traces.
package telemetry

import (
	"context"
	"os"
	"sync"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	loggerProvider *sdklog.LoggerProvider
	loggerOnce     sync.Once
)

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	ServiceName string
	Endpoint    string // OTLP endpoint, empty for stdout
}

// InitLogging initializes OpenTelemetry logging with OTLP or stdout exporter.
func InitLogging(ctx context.Context, cfg LoggingConfig) error {
	var err error
	loggerOnce.Do(func() {
		err = initLoggingInternal(ctx, cfg)
	})
	return err
}

func initLoggingInternal(ctx context.Context, cfg LoggingConfig) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return err
	}

	var exporter sdklog.Exporter
	if cfg.Endpoint != "" {
		// Use OTLP HTTP exporter for production
		exporter, err = otlploghttp.New(ctx,
			otlploghttp.WithEndpoint(cfg.Endpoint),
		)
	} else {
		// Use stdout exporter for development/testing
		exporter, err = stdoutlog.New(
			stdoutlog.WithWriter(os.Stdout),
		)
	}
	if err != nil {
		return err
	}

	loggerProvider = sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	global.SetLoggerProvider(loggerProvider)

	return nil
}

// ShutdownLogging gracefully shuts down the logger provider.
func ShutdownLogging(ctx context.Context) error {
	if loggerProvider != nil {
		return loggerProvider.Shutdown(ctx)
	}
	return nil
}

// GetLoggerProvider returns the configured logger provider.
func GetLoggerProvider() *sdklog.LoggerProvider {
	return loggerProvider
}
