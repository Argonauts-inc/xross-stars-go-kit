package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	globalotel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Provider struct {
	shutdown func(context.Context) error
}

func (p *Provider) Shutdown(ctx context.Context) error {
	if p == nil || p.shutdown == nil {
		return nil
	}
	return p.shutdown(ctx)
}

func Init(ctx context.Context, cfg Config) (*Provider, error) {
	cfg = normalizeConfig(cfg)
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}
	if cfg.TraceSampleRate < 0 || cfg.TraceSampleRate > 1 {
		return nil, fmt.Errorf("trace sample rate must be between 0 and 1")
	}

	globalotel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		xray.Propagator{},
	))

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
		semconv.DeploymentEnvironment(cfg.Environment),
		semconv.ServiceVersion(cfg.ServiceVersion),
		attribute.String("service.namespace", cfg.ServiceNamespace),
	)

	options := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.TraceSampleRate))),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
	}

	if cfg.TraceExporter {
		exporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			return nil, err
		}
		options = append(options, sdktrace.WithBatcher(
			exporter,
			sdktrace.WithBatchTimeout(cfg.TraceExportPeriod),
		))
	}

	provider := sdktrace.NewTracerProvider(options...)
	globalotel.SetTracerProvider(provider)

	return &Provider{shutdown: provider.Shutdown}, nil
}
