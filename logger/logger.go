package logger

import (
	"context"
	"io"
	"log/slog"

	"github.com/Argonauts-inc/xross-stars-go-kit/requestid"
	"go.opentelemetry.io/otel/trace"
)

const (
	TraceIDLogKey = "trace_id"
	SpanIDLogKey  = "span_id"
)

func NewJSONLogger(w io.Writer, opts *slog.HandlerOptions) *slog.Logger {
	return slog.New(NewJSONHandler(w, opts))
}

func NewJSONHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return NewHandler(slog.NewJSONHandler(w, opts))
}

func NewHandler(next slog.Handler) slog.Handler {
	return requestid.NewLogHandler(NewTraceHandler(next))
}

type traceHandler struct {
	next slog.Handler
}

func NewTraceHandler(next slog.Handler) slog.Handler {
	return &traceHandler{next: next}
}

func (h *traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *traceHandler) Handle(ctx context.Context, record slog.Record) error {
	if traceID, spanID, ok := TraceFields(ctx); ok {
		record.AddAttrs(
			slog.String(TraceIDLogKey, traceID),
			slog.String(SpanIDLogKey, spanID),
		)
	}

	return h.next.Handle(ctx, record)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{next: h.next.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{next: h.next.WithGroup(name)}
}

func TraceFields(ctx context.Context) (traceID string, spanID string, ok bool) {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return "", "", false
	}

	return spanContext.TraceID().String(), spanContext.SpanID().String(), true
}
