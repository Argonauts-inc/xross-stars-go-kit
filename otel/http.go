package otel

import (
	"net/http"

	globalotel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const DefaultHTTPServerTracerName = "xross-stars/http-server"

type RouteNameFunc func(method, path string) (pattern string, ok bool)

type HTTPTraceConfig struct {
	TracerName string
	RouteName  RouteNameFunc
}

func HTTPTraceContextMiddleware(next http.Handler, routeName RouteNameFunc) http.Handler {
	return HTTPTraceMiddleware(next, HTTPTraceConfig{RouteName: routeName})
}

func HTTPTraceMiddleware(next http.Handler, cfg HTTPTraceConfig) http.Handler {
	if cfg.TracerName == "" {
		cfg.TracerName = DefaultHTTPServerTracerName
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := globalotel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		spanName := r.Method
		attrs := []attribute.KeyValue{
			attribute.String("http.request.method", r.Method),
			attribute.String("url.path", r.URL.Path),
			attribute.String("url.query", r.URL.RawQuery),
			attribute.String("user_agent.original", r.UserAgent()),
		}
		if cfg.RouteName != nil {
			if pattern, ok := cfg.RouteName(r.Method, r.URL.Path); ok {
				spanName = r.Method + " " + pattern
				attrs = append(attrs, attribute.String("http.route", pattern))
			}
		}

		ctx, span := globalotel.Tracer(cfg.TracerName).Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r.WithContext(ctx))

		span.SetAttributes(attribute.Int("http.response.status_code", recorder.statusCode))
		if recorder.statusCode >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(recorder.statusCode))
		}
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}

	r.statusCode = statusCode
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}

	return r.ResponseWriter.Write(b)
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}
