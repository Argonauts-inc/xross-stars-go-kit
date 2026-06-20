package otel

import (
	"context"
	"strings"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/propagation"
)

const XRayTraceHeader = "AWSTraceHeader"

func InjectXRayTraceHeader(ctx context.Context) string {
	carrier := propagation.MapCarrier{}
	(xray.Propagator{}).Inject(ctx, carrier)

	for _, field := range (xray.Propagator{}).Fields() {
		if value := carrier.Get(field); value != "" {
			return value
		}
	}

	return ""
}

func ExtractXRayTraceHeader(ctx context.Context, traceHeader string) context.Context {
	if strings.TrimSpace(traceHeader) == "" {
		return ctx
	}

	carrier := propagation.MapCarrier{}
	for _, field := range (xray.Propagator{}).Fields() {
		carrier.Set(field, traceHeader)
	}

	return (xray.Propagator{}).Extract(ctx, carrier)
}
