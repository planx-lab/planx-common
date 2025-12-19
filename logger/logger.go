// Package logger provides a zerolog-based logging implementation for Planx.
// All Planx modules MUST use this logger for consistent logging.
// Integrates with OpenTelemetry for trace context and log export.
// Engine-side logging utilities.
// MUST NOT be imported by SDK or plugins.

package logger

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var (
	globalLogger zerolog.Logger
	once         sync.Once
)

// Config holds logger configuration.
type Config struct {
	Level       string // debug, info, warn, error
	Pretty      bool   // human-readable output (for development)
	Output      io.Writer
	ServiceName string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Level:       "info",
		Pretty:      false,
		Output:      os.Stdout,
		ServiceName: "planx",
	}
}

// Init initializes the global logger with the given configuration.
func Init(cfg Config) {
	once.Do(func() {
		level, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			level = zerolog.InfoLevel
		}

		zerolog.SetGlobalLevel(level)
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var output io.Writer = cfg.Output
		if cfg.Pretty {
			output = zerolog.ConsoleWriter{
				Out:        cfg.Output,
				TimeFormat: "15:04:05.000",
			}
		}

		globalLogger = zerolog.New(output).
			With().
			Timestamp().
			Str("service", cfg.ServiceName).
			Logger()
	})
}

// Get returns the global logger.
func Get() *zerolog.Logger {
	if globalLogger.GetLevel() == zerolog.Disabled {
		Init(DefaultConfig())
	}
	return &globalLogger
}

// WithContext returns a logger with OpenTelemetry trace context fields.
// Automatically extracts trace_id and span_id from the context if present.
// This enables log correlation with distributed traces.
func WithContext(ctx context.Context) *zerolog.Logger {
	l := Get().With().Logger()

	// Extract OpenTelemetry trace context
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		l = l.With().Str("trace_id", span.SpanContext().TraceID().String()).Logger()
	}
	if span.SpanContext().HasSpanID() {
		l = l.With().Str("span_id", span.SpanContext().SpanID().String()).Logger()
	}

	return &l
}

// ContextWithTrace returns a context with trace information embedded (legacy support).
// Prefer using OpenTelemetry context propagation instead.
func ContextWithTrace(ctx context.Context, traceID, spanID string) context.Context {
	ctx = context.WithValue(ctx, traceIDKey, traceID)
	ctx = context.WithValue(ctx, spanIDKey, spanID)
	return ctx
}

// Context keys for tracing (legacy)
type contextKey string

const (
	traceIDKey contextKey = "trace_id"
	spanIDKey  contextKey = "span_id"
)

// Helper functions for direct usage (without context)

// Debug logs at debug level.
func Debug() *zerolog.Event {
	return Get().Debug()
}

// Info logs at info level.
func Info() *zerolog.Event {
	return Get().Info()
}

// Warn logs at warn level.
func Warn() *zerolog.Event {
	return Get().Warn()
}

// Error logs at error level.
func Error() *zerolog.Event {
	return Get().Error()
}

// Fatal logs at fatal level and exits.
func Fatal() *zerolog.Event {
	return Get().Fatal()
}

// Context-aware helper functions
// These automatically inject trace_id and span_id from OTel context.

// DebugCtx logs at debug level with trace context.
func DebugCtx(ctx context.Context) *zerolog.Event {
	return WithContext(ctx).Debug()
}

// InfoCtx logs at info level with trace context.
func InfoCtx(ctx context.Context) *zerolog.Event {
	return WithContext(ctx).Info()
}

// WarnCtx logs at warn level with trace context.
func WarnCtx(ctx context.Context) *zerolog.Event {
	return WithContext(ctx).Warn()
}

// ErrorCtx logs at error level with trace context.
func ErrorCtx(ctx context.Context) *zerolog.Event {
	return WithContext(ctx).Error()
}

// AddSpanEvent adds a log message as a span event for OTel correlation.
// This bridges zerolog logs to OTel spans.
func AddSpanEvent(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(msg)
	}
}

// AddSpanEventWithAttrs adds a log message with attributes as a span event.
func AddSpanEventWithAttrs(ctx context.Context, msg string, attrs map[string]string) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		// Note: For full OTel log export, use OTLP log exporter
		// This is a simple bridge that adds log events to spans
		span.AddEvent(msg)
	}
}
