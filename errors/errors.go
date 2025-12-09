// Package errors provides error types with stack traces for Planx.
package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Error represents an error with a stack trace and optional cause.
type Error struct {
	Message string
	Cause   error
	Stack   []uintptr
}

// New creates a new error with a stack trace.
func New(message string) *Error {
	return &Error{
		Message: message,
		Stack:   captureStack(2),
	}
}

// Wrap wraps an existing error with additional context and a stack trace.
func Wrap(err error, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Message: message,
		Cause:   err,
		Stack:   captureStack(2),
	}
}

// Wrapf wraps an existing error with formatted context.
func Wrapf(err error, format string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
		Stack:   captureStack(2),
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	return e.Cause
}

// StackTrace returns a formatted stack trace.
func (e *Error) StackTrace() string {
	var sb strings.Builder
	frames := runtime.CallersFrames(e.Stack)
	for {
		frame, more := frames.Next()
		if frame.Function == "" {
			break
		}
		sb.WriteString(fmt.Sprintf("  %s\n    %s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}

func captureStack(skip int) []uintptr {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skip+1, pcs)
	return pcs[:n]
}

// Error types for categorization

// ConfigError represents a configuration error (fatal on CreateSession).
type ConfigError struct {
	*Error
}

// NewConfigError creates a new configuration error.
func NewConfigError(message string) *ConfigError {
	return &ConfigError{Error: New(message)}
}

// StreamError represents a stream error (terminate session).
type StreamError struct {
	*Error
}

// NewStreamError creates a new stream error.
func NewStreamError(message string) *StreamError {
	return &StreamError{Error: New(message)}
}

// BatchError represents a batch-level error (partial failure allowed).
type BatchError struct {
	*Error
	FailedIndices []int
}

// NewBatchError creates a new batch error with failed record indices.
func NewBatchError(message string, failedIndices []int) *BatchError {
	return &BatchError{
		Error:         New(message),
		FailedIndices: failedIndices,
	}
}

// TransportError represents a transport error (retry connection).
type TransportError struct {
	*Error
	Retryable bool
}

// NewTransportError creates a new transport error.
func NewTransportError(message string, retryable bool) *TransportError {
	return &TransportError{
		Error:     New(message),
		Retryable: retryable,
	}
}
