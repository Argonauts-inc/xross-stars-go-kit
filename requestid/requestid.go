package requestid

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

const (
	HeaderName = "X-Request-Id"
	LogKey     = "request_id"
	MaxLength  = 128
)

type contextKey struct{}

func New() string {
	return uuid.NewString()
}

func Normalize(id string) string {
	id = strings.TrimSpace(id)
	if id == "" || len(id) > MaxLength {
		return ""
	}

	for _, r := range id {
		if r < 0x21 || r == 0x7f {
			return ""
		}
	}

	return id
}

func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKey{}).(string)
	if !ok || id == "" {
		return "", false
	}

	return id, true
}

func WithContext(ctx context.Context, id string) context.Context {
	if id = Normalize(id); id == "" {
		return ctx
	}

	return context.WithValue(ctx, contextKey{}, id)
}

func Ensure(ctx context.Context, id string) (context.Context, string) {
	if id = Normalize(id); id == "" {
		id = New()
	}

	return WithContext(ctx, id), id
}

type logHandler struct {
	next slog.Handler
}

func NewLogHandler(next slog.Handler) slog.Handler {
	return &logHandler{next: next}
}

func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *logHandler) Handle(ctx context.Context, record slog.Record) error {
	if id, ok := FromContext(ctx); ok {
		record.AddAttrs(slog.String(LogKey, id))
	}

	return h.next.Handle(ctx, record)
}

func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{next: h.next.WithAttrs(attrs)}
}

func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{next: h.next.WithGroup(name)}
}
