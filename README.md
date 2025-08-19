
# bind

`bind` is a Go library for binding [OpenTelemetry] synchronous metric instruments to attributes. This can greatly improver performance for instrumentation when measurements are only recorded for known static attribute sets.

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

counter = bind.Float64Counter(counter, attribute.Bool("admin", true))

// Measure with {"user": "Alice", "admin": true}
counter.Add(ctx, -1.0)

// Unwrap to get the underlying instrument (i.e. myCounter) and attributes (i.e. {"user": "Alice", "admin": true})
base, attrs := bind.Unwrap(counter)
```

See [GoDoc](https://pkg.go.dev/github.com/MrAlias/bind) for full API documentation and examples.

## Contributing

Contributions are welcome! Please open issues or pull requests on [GitHub](https://github.com/MrAlias/bind).

## License

Licensed under the [Apache 2.0 License](LICENSE).

[OpenTelemetry]: opentelemetry.io
