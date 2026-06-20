package requestid

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   string
		want string
	}{
		{name: "trims spaces", id: " req-123 ", want: "req-123"},
		{name: "empty", id: " ", want: ""},
		{name: "too long", id: strings.Repeat("a", MaxLength+1), want: ""},
		{name: "control character", id: "bad\nid", want: ""},
		{name: "delete character", id: "bad" + string(rune(0x7f)), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Normalize(tt.id); got != tt.want {
				t.Fatalf("Normalize(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestWithContext(t *testing.T) {
	t.Parallel()

	ctx := WithContext(context.Background(), " req-123 ")
	got, ok := FromContext(ctx)
	if !ok {
		t.Fatal("request id missing from context")
	}
	if got != "req-123" {
		t.Fatalf("request id = %q, want req-123", got)
	}
}

func TestWithContextIgnoresInvalidID(t *testing.T) {
	t.Parallel()

	ctx := WithContext(context.Background(), "bad\nid")
	if id, ok := FromContext(ctx); ok {
		t.Fatalf("request id = %q, want missing", id)
	}
}

func TestEnsureGeneratesID(t *testing.T) {
	t.Parallel()

	ctx, id := Ensure(context.Background(), "")
	if id == "" {
		t.Fatal("generated id is empty")
	}

	got, ok := FromContext(ctx)
	if !ok {
		t.Fatal("generated id missing from context")
	}
	if got != id {
		t.Fatalf("context id = %q, want %q", got, id)
	}
}

func TestLogHandlerAddsRequestID(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(NewLogHandler(slog.NewJSONHandler(&buf, nil)))
	ctx := WithContext(context.Background(), "req-123")

	logger.InfoContext(ctx, "hello")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal log: %v", err)
	}
	if got := entry[LogKey]; got != "req-123" {
		t.Fatalf("%s = %v, want req-123", LogKey, got)
	}
}
