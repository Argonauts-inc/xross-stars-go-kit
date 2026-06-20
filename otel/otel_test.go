package otel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	globalotel "go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestTraceExporterEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		endpoint       string
		tracesEndpoint string
		want           bool
	}{
		{name: "disabled", want: false},
		{name: "endpoint", endpoint: "localhost:4317", want: true},
		{name: "traces endpoint", tracesEndpoint: "localhost:4318", want: true},
		{name: "trim spaces", endpoint: "  ", tracesEndpoint: " endpoint ", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := TraceExporterEnabled(tt.endpoint, tt.tracesEndpoint); got != tt.want {
				t.Fatalf("TraceExporterEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitRequiresServiceName(t *testing.T) {
	t.Parallel()

	if _, err := Init(context.Background(), Config{}); err == nil {
		t.Fatal("Init() error = nil, want error")
	}
}

func TestHTTPTraceContextMiddlewareUsesRouteTemplate(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	originalProvider := globalotel.GetTracerProvider()
	globalotel.SetTracerProvider(provider)
	t.Cleanup(func() {
		globalotel.SetTracerProvider(originalProvider)
		_ = provider.Shutdown(context.Background())
	})

	handler := HTTPTraceContextMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}), func(method, path string) (string, bool) {
		if method == http.MethodGet && path == "/v1/cards/123" {
			return "/v1/cards/{id}", true
		}
		return "", false
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/cards/123?foo=bar", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("ended spans = %d, want 1", len(spans))
	}
	span := spans[0]
	if span.Name() != "GET /v1/cards/{id}" {
		t.Fatalf("span name = %q, want GET /v1/cards/{id}", span.Name())
	}
	if span.SpanKind() != trace.SpanKindServer {
		t.Fatalf("span kind = %v, want server", span.SpanKind())
	}
}

func TestInjectAndExtractXRayTraceHeader(t *testing.T) {
	t.Parallel()

	provider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	ctx, span := provider.Tracer("test").Start(context.Background(), "test")
	defer span.End()

	traceHeader := InjectXRayTraceHeader(ctx)
	if traceHeader == "" {
		t.Fatal("trace header is empty")
	}

	extracted := ExtractXRayTraceHeader(context.Background(), traceHeader)
	spanContext := trace.SpanContextFromContext(extracted)
	if !spanContext.IsValid() {
		t.Fatal("extracted span context is invalid")
	}
}
