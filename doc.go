/*
Package bind provides functions for binding OpenTelemetry synchronous metric
instruments to attributes. This can greatly improve performance for
instrumentation when measurements are only recorded for known static attribute
sets.

The package allows you to pre-bind attributes to metric instruments, so that
all subsequent measurements automatically include those attributes without
needing to specify them each time.

Example usage:

	// Bind a user attribute to a counter.
	user := attribute.String("user", "Alice"))
	counter := bind.Float64Counter(myCounter, user)

	// All measurements will include {"user": "Alice"}
	counter.Add(ctx, 1.0)

	// Additional attributes can be added at measurement time.
	counter.Add(ctx, 2.0, metric.WithAttributes(attribute.Int("id", 1)))

The package supports all OpenTelemetry synchronous metric instrument types:

  - Int64Counter
  - Int64UpDownCounter
  - Int64Histogram
  - Int64Gauge
  - Float64Counter
  - Float64UpDownCounter
  - Float64Histogram
  - Float64Gauge.

Bound instruments can be further bound with additional attributes, or the
original instrument and attributes can be retrieved using [Unwrap].
*/
package bind
