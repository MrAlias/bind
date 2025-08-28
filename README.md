# bind

[![CI](https://github.com/MrAlias/bind/actions/workflows/ci.yml/badge.svg)](https://github.com/MrAlias/bind/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/MrAlias/bind/branch/main/graph/badge.svg)](https://codecov.io/gh/MrAlias/bind)
[![Go Report Card](https://goreportcard.com/badge/github.com/MrAlias/bind)](https://goreportcard.com/report/github.com/MrAlias/bind)
[![Go Reference](https://pkg.go.dev/badge/github.com/MrAlias/bind.svg)](https://pkg.go.dev/github.com/MrAlias/bind)
[![GitHub release](https://img.shields.io/github/release/MrAlias/bind.svg)](https://github.com/MrAlias/bind/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`bind` is a Go library for binding [OpenTelemetry] synchronous metric instruments to attributes.
This can greatly improve performance for instrumentation when measurements are only recorded for known static attribute sets.

## Installation

```sh
go get github.com/MrAlias/bind
```

## Usage

```go
import (
 "github.com/MrAlias/bind"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/metric"
)

// Example: Bind attributes to a Float64Counter
counter := bind.Float64Counter(myCounter, attribute.String("user", "Alice"))

// Measured with {"user": "Alice"}
counter.Add(ctx, 1.0)

// Measured with {"user": "Alice", "id: 1}
counter.Add(ctx, 2.0, metric.WithAttributes(attribute.Int("id", 1)))

// Bind additional attributes to the counter.
counter = bind.Float64Counter(counter, attribute.Bool("admin", true))

// Measure with {"user": "Alice", "admin": true}
counter.Add(ctx, -1.0)

// Unwrap to get the underlying instrument (i.e. myCounter) and attributes
// (i.e. {"user": "Alice", "admin": true})
base, attrs := bind.Unwrap(counter)
```

See [GoDoc] for full API documentation and examples.

## Contributing

Contributions are welcome!
Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

Licensed under the [Apache 2.0 License](LICENSE).

[OpenTelemetry]: opentelemetry.io
[GoDoc]: https://pkg.go.dev/github.com/MrAlias/bind
