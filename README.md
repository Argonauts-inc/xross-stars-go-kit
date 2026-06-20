# xross-stars-go-kit

Go Utilities Kit for Xross Stars.

This repository provides reusable Go helpers for Xross Stars services. It is
intended to be safe to operate as OSS: only public, generic implementation
patterns belong here. Service-specific configuration, secrets, internal URLs,
AWS account details, production payloads, logs, traces, and customer data must
stay in the consuming repositories or runtime configuration.

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

The module currently targets Go `1.24.1` so it can be consumed by existing
Xross Stars Lambda functions that still build with Go 1.24.x.

OpenTelemetry dependencies are pinned to the 1.37.x generation in this kit. Go's
minimal version selection means services that already require newer compatible
OpenTelemetry modules, such as Xross Stars API services using 1.44.x, will keep
using the newer version.

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

This is a Go library module. Tags are the module versions consumed by downstream
projects.

To release a new version:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Pushing a `v*` tag runs the release workflow. GoReleaser creates the GitHub
Release and changelog for the tag. Because this repository does not ship a CLI
binary, GoReleaser is configured to skip binary builds.
