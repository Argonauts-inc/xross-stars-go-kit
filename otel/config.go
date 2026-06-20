package otel

import (
	"strings"
	"time"
)

const (
	DefaultServiceNamespace = "xross-stars"
	DefaultTraceExportDelay = 5 * time.Second
)

type Config struct {
	ServiceName       string
	Environment       string
	ServiceVersion    string
	ServiceNamespace  string
	TraceSampleRate   float64
	TraceExporter     bool
	TraceExportPeriod time.Duration
}

func TraceExporterEnabled(endpoint, tracesEndpoint string) bool {
	return strings.TrimSpace(endpoint) != "" || strings.TrimSpace(tracesEndpoint) != ""
}

func normalizeConfig(cfg Config) Config {
	if cfg.ServiceNamespace == "" {
		cfg.ServiceNamespace = DefaultServiceNamespace
	}
	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = "unknown"
	}
	if cfg.TraceExportPeriod == 0 {
		cfg.TraceExportPeriod = DefaultTraceExportDelay
	}

	return cfg
}
