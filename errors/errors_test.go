package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	e := New("something failed")
	if e.Message != "something failed" {
		t.Fatalf("message: got %q", e.Message)
	}
	if e.Cause != nil {
		t.Fatalf("cause: got %v, want nil", e.Cause)
	}
	if len(e.Stack) == 0 {
		t.Fatal("stack should not be empty")
	}
}

func TestNew_ErrorString(t *testing.T) {
	e := New("msg")
	if e.Error() != "msg" {
		t.Fatalf("got %q", e.Error())
	}
}

func TestNew_WithCause(t *testing.T) {
	inner := fmt.Errorf("inner")
	e := New("outer")
	e.Cause = inner
	got := e.Error()
	if !strings.Contains(got, "outer") || !strings.Contains(got, "inner") {
		t.Fatalf("got %q", got)
	}
}

func TestWrap(t *testing.T) {
	inner := fmt.Errorf("root")
	e := Wrap(inner, "wrapped")
	if e.Message != "wrapped" {
		t.Fatalf("message: got %q", e.Message)
	}
	if e.Cause != inner {
		t.Fatalf("cause mismatch")
	}
	if len(e.Stack) == 0 {
		t.Fatal("stack should not be empty")
	}
}

func TestWrap_Nil(t *testing.T) {
	e := Wrap(nil, "msg")
	if e != nil {
		t.Fatalf("got %v, want nil", e)
	}
}

func TestWrapf(t *testing.T) {
	inner := fmt.Errorf("root")
	e := Wrapf(inner, "code %d", 42)
	if e.Message != "code 42" {
		t.Fatalf("message: got %q", e.Message)
	}
	if e.Cause != inner {
		t.Fatal("cause mismatch")
	}
}

func TestWrapf_Nil(t *testing.T) {
	e := Wrapf(nil, "msg")
	if e != nil {
		t.Fatalf("got %v, want nil", e)
	}
}

func TestError_Unwrap(t *testing.T) {
	inner := fmt.Errorf("root")
	e := Wrap(inner, "outer")
	if !errors.Is(e, inner) {
		t.Fatal("errors.Is should match inner")
	}
}

func TestNil_Error(t *testing.T) {
	var e *Error
	if e.Error() != "" {
		t.Fatalf("got %q, want empty", e.Error())
	}
}

func TestNil_Unwrap(t *testing.T) {
	var e *Error
	if e.Unwrap() != nil {
		t.Fatal("expected nil")
	}
}

func TestNil_StackTrace(t *testing.T) {
	var e *Error
	if e.StackTrace() != "" {
		t.Fatalf("got %q, want empty", e.StackTrace())
	}
}

func TestStackTrace_ContainsFunction(t *testing.T) {
	e := New("trace me")
	st := e.StackTrace()
	if st == "" {
		t.Fatal("stack trace empty")
	}
	if !strings.Contains(st, "TestStackTrace_ContainsFunction") {
		t.Fatalf("stack trace should contain test function name, got:\n%s", st)
	}
}

func TestNewConfigError(t *testing.T) {
	e := NewConfigError("bad config")
	if e.Error.Message != "bad config" {
		t.Fatalf("message: got %q", e.Error.Message)
	}
	if e.Error == nil {
		t.Fatal("embedded Error should not be nil")
	}
}

func TestNewStreamError(t *testing.T) {
	e := NewStreamError("stream broke")
	if e.Error.Message != "stream broke" {
		t.Fatalf("message: got %q", e.Error.Message)
	}
}

func TestNewBatchError(t *testing.T) {
	indices := []int{2, 5, 7}
	e := NewBatchError("partial fail", indices)
	if e.Error.Message != "partial fail" {
		t.Fatalf("message: got %q", e.Error.Message)
	}
	if len(e.FailedIndices) != 3 || e.FailedIndices[0] != 2 {
		t.Fatalf("indices: got %v", e.FailedIndices)
	}
}

func TestNewBatchError_EmptyIndices(t *testing.T) {
	e := NewBatchError("empty", nil)
	if e.FailedIndices != nil {
		t.Fatalf("got %v, want nil", e.FailedIndices)
	}
}

func TestNewTransportError(t *testing.T) {
	e := NewTransportError("timeout", true)
	if e.Error.Message != "timeout" {
		t.Fatalf("message: got %q", e.Error.Message)
	}
	if !e.Retryable {
		t.Fatal("should be retryable")
	}
}

func TestNewTransportError_NotRetryable(t *testing.T) {
	e := NewTransportError("refused", false)
	if e.Retryable {
		t.Fatal("should not be retryable")
	}
}

func TestWrappedError_Unwrap(t *testing.T) {
	inner := fmt.Errorf("base")
	e := Wrap(inner, "layer1")
	e2 := Wrap(e, "layer2")
	if !errors.Is(e2, inner) {
		t.Fatal("should unwrap to base")
	}
	if !errors.Is(e2, e) {
		t.Fatal("should unwrap to layer1")
	}
}

func TestConfigError_CallsEmbedded(t *testing.T) {
	e := NewConfigError("cfg")
	got := e.Error.Error()
	if got != "cfg" {
		t.Fatalf("got %q", got)
	}
}

func TestBatchError_CallsEmbedded(t *testing.T) {
	e := NewBatchError("batch", []int{1})
	got := e.Error.Error()
	if got != "batch" {
		t.Fatalf("got %q", got)
	}
}

func TestTransportError_CallsEmbedded(t *testing.T) {
	e := NewTransportError("trans", false)
	got := e.Error.Error()
	if got != "trans" {
		t.Fatalf("got %q", got)
	}
}
