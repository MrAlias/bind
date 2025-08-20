package bind_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/MrAlias/bind"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// Example demonstrates how to bind attributes to a metric instrument and use
// it.
func Example() {
	ctx := context.Background()
	// As an example, we assume the user has already create a Float64Counter
	// instrument they want to bind to attributes.
	var myCounter metric.Float64Counter = noop.Float64Counter{}

	// Bind the instrument to attributes.
	counter := bind.Float64Counter(
		myCounter,
		attribute.String("user", "Alice"),
		attribute.Int("id", 73),
	)

	// Use the bound instrument to record a measurement. This will measure 1.0
	// for the {user: "Alice", "id": 73} attributes.
	counter.Add(ctx, 1.0)

	// Dynmic attributes can be added at any time the same as with any metric
	// instrument. This will measure 2.0 for the {user: "Alice", "id": 73,
	// score: 98.3} attributes.
	score := attribute.Float64("score", 98.3)
	counter.Add(ctx, 2.0, metric.WithAttributes(score))

	// A new bound instrument can be created with additional attributes.
	counter = bind.Float64Counter(counter, attribute.Bool("admin", true))

	// This will measure 0.1 for the {user: "Alice", "id": 73, "admin": true}
	// attributes.
	counter.Add(ctx, 0.1)

	// Call bind.Unwrap to get the original instrument and the
	// attributes that were bound to it.
	original, attrs := bind.Unwrap(counter)

	var buf []string
	for _, attr := range attrs.ToSlice() {
		str := fmt.Sprintf("%s = %v", attr.Key, attr.Value.AsInterface())
		buf = append(buf, str)
	}
	fmt.Printf("Original Instrument: %T\n", original)
	fmt.Printf("Bound Attributes: {%s}\n", strings.Join(buf, ", "))
	// Output:
	//
	// Original Instrument: noop.Float64Counter
	// Bound Attributes: {admin = true, id = 73, user = Alice}
}
