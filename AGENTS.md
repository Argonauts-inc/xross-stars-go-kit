# AGENTS.md

This file provides guidance to coding agents when working with this repository.

## Project Overview

This repository is "xross-stars-go-kit", a Go utilities kit for Xross Stars.
It contains reusable Go helpers that can be used by Xross Stars services without
copying observability, logging, request ID, AWS SQS propagation, or database
instrumentation code between repositories.

This repository is intended to be operated as OSS. Only public, reusable
implementation details belong here.

## Development Commands

```bash
go test ./...
go vet ./...
gofmt -w .
go mod tidy
goreleaser check
```

## Package Layout

| Package | Purpose |
|---------|---------|
| `requestid` | Request ID normalization, context helpers, and slog handler |
| `logger` | slog helpers that attach request_id, trace_id, and span_id |
| `otel` | OpenTelemetry provider setup, HTTP tracing middleware, and X-Ray propagation helpers |
| `awssqs` | AWS SQS request_id and X-Ray trace propagation helpers |
| `pgxotel` | pgx / otelpgx helper logic such as low-cardinality SQL span names |

## OSS and Confidentiality Rules

- Do not commit secrets, credentials, tokens, private keys, or `.env` values.
- Do not hard-code AWS account IDs, ARNs, bucket names, queue URLs, database URLs,
  production endpoints, internal hostnames, or customer/user data.
- Do not copy private business payloads, logs, traces, incident data, or
  non-public operational details from Xross Stars.
- Keep service-specific values configurable from the caller. This kit should
  expose options and helpers, not Xross Stars runtime configuration.
- Prefer generic package names and public APIs that make sense outside a single
  service implementation.
- Keep dependencies small and purposeful. Avoid adding cloud/provider libraries
  to a package unless that package is explicitly scoped to that provider.

## Engineering Guidelines

- Preserve compatibility for downstream Xross Stars services when changing
  exported APIs.
- Add or update tests for every exported helper and behavior change.
- Keep span and metric names low-cardinality. Put detailed values in attributes
  only when they are safe to expose and useful for debugging.
- Logging helpers must not log sensitive headers, request bodies, credentials, or
  business payloads by default.
- Public APIs should use standard library types where practical.

## Release Notes

This repository is consumed as a normal external Go module:

```go
import "github.com/Argonauts-inc/xross-stars-go-kit/otel"
```

Tag releases when downstream Xross Stars repositories need stable versions.
Pushing a `v*` tag runs the GitHub Actions release workflow and GoReleaser
creates the GitHub Release and changelog. This repository is a library, so
GoReleaser is configured to skip binary builds.
