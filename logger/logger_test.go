package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/Argonauts-inc/xross-stars-go-kit/requestid"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestNewJSONLoggerAddsRequestAndTraceFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	ctx, span := provider.Tracer("test").Start(context.Background(), "test-span")
	defer span.End()
	ctx = requestid.WithContext(ctx, "req-123")

	NewJSONLogger(&buf, nil).InfoContext(ctx, "hello")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal log: %v", err)
	}
	if got := entry[requestid.LogKey]; got != "req-123" {
		t.Fatalf("%s = %v, want req-123", requestid.LogKey, got)
	}
	if got := entry[TraceIDLogKey]; got == "" {
		t.Fatalf("%s is empty", TraceIDLogKey)
	}
	if got := entry[SpanIDLogKey]; got == "" {
		t.Fatalf("%s is empty", SpanIDLogKey)
	}
}

func TestTraceFieldsReturnsFalseWithoutSpan(t *testing.T) {
	t.Parallel()

	if traceID, spanID, ok := TraceFields(context.Background()); ok {
		t.Fatalf("TraceFields ok with trace_id=%q span_id=%q, want false", traceID, spanID)
	}
}
