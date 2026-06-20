# xross-stars-go-kit

Go Utilities Kit for Xross Stars.

This module provides small, focused Go packages for request correlation,
structured logging, OpenTelemetry setup, AWS SQS trace propagation, and pgx
observability. Each package is designed to be imported independently and
configured by the consuming application.

## Packages

| Package | Purpose |
|---------|---------|
| [`requestid`](./requestid) | Request ID normalization, context helpers, and `slog` handler |
| [`logger`](./logger) | JSON `slog` helpers that attach `request_id`, `trace_id`, and `span_id` |
| [`otel`](./otel) | OpenTelemetry provider setup, HTTP tracing middleware, and AWS X-Ray trace header helpers |
| [`awssqs`](./awssqs) | AWS SQS request ID and X-Ray trace propagation helpers |
| [`pgxotel`](./pgxotel) | Helpers for low-cardinality SQL span names with pgx / otelpgx |

## Installation

Use this repository as a normal Go module:

```bash
go get github.com/Argonauts-inc/xross-stars-go-kit@latest
```

Import only the package you need:

```go
import (
	"github.com/Argonauts-inc/xross-stars-go-kit/logger"
	"github.com/Argonauts-inc/xross-stars-go-kit/otel"
	"github.com/Argonauts-inc/xross-stars-go-kit/requestid"
)
```

## Compatibility

The module currently targets Go `1.24.1` for compatibility with Go 1.24+
projects.

OpenTelemetry dependencies are pinned to the 1.37.x generation in this kit. Go's
minimal version selection means services that already require newer compatible
OpenTelemetry modules will keep using the newer version.

## Development

```bash
gofmt -w .
go vet ./...
go test ./...
goreleaser check
```

## CI

GitHub Actions runs formatting checks, `go vet`, `go test`, and GoReleaser
configuration validation on pull requests and pushes to `main`.

## Release

This is a Go library module. Tags are the module versions consumed by Go
projects.

To release a new version:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Pushing a `v*` tag runs the release workflow. GoReleaser creates the GitHub
Release and changelog for the tag. Because this repository does not ship a CLI
binary, GoReleaser is configured to skip binary builds.
