package bind_test

import (
	"context"
	"testing"

	"github.com/MrAlias/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/noop"
)

type mockFloat64UpDownCounter struct {
	embedded.Float64UpDownCounter

	incr *float64
	opts []metric.AddOption
}

func (m *mockFloat64UpDownCounter) Add(ctx context.Context, incr float64, opts ...metric.AddOption) {
	m.incr = &incr
	m.opts = opts
}

func (m *mockFloat64UpDownCounter) Instrument() metric.Float64UpDownCounter {
	return m
}

func (m *mockFloat64UpDownCounter) Recorded() (*float64, []attribute.KeyValue) {
	set := metric.NewAddConfig(m.opts).Attributes()
	return m.incr, set.ToSlice()
}

func TestFloat64UpDownCounterEmptyAttrs(t *testing.T) {
	mock := &mockFloat64UpDownCounter{}
	got := bind.Float64UpDownCounter(mock)
	assert.Same(t, mock, got, "bound should be the same as the input")
}

func TestFloat64UpDownCounter(t *testing.T) {
	tests := []TestCase[float64]{
		{"BoundOnly", 10.0, nil},
		{"AddAttr", -5.0, []attribute.KeyValue{adminTrue}},
	}

	for _, test := range tests {
		t.Run(test.Name, Run(
			&mockFloat64UpDownCounter{},
			bind.Float64UpDownCounter,
			measFloat64UpDownCounter,
			test,
		))
	}

	t.Run(
		"NoSideEffets",
		testNoSideEffects(bind.Float64UpDownCounter, &mockFloat64UpDownCounter{}),
	)
}

func measFloat64UpDownCounter(i metric.Float64UpDownCounter, ctx context.Context, incr float64, attr []attribute.KeyValue) {
	var opts []metric.AddOption
	if attr != nil {
		opts = append(opts, metric.WithAttributes(attr...))
	}
	i.Add(ctx, incr, opts...)
}

func TestFloat64UpDownCounterUnwrap(t *testing.T) {
	mock := &mockFloat64UpDownCounter{}
	bound := bind.Float64UpDownCounter(mock, userAlice)
	require.NotNil(t, bound, "bound should not be nil")
	bound = bind.Float64UpDownCounter(bound, userID)
	require.NotNil(t, bound, "bound should not be nil")

	val, set := bind.Unwrap(bound)

	assert.Same(t, mock, val, "unwrapped value should match mock")

	want := []attribute.KeyValue{userAlice, userID}
	assert.ElementsMatch(t, want, set.ToSlice(), "unwrapped attributes")
}

func BenchmarkFloat64UpDownCounterAdd(b *testing.B) {
	ctx := context.Background()
	base := noop.Float64UpDownCounter{}

	b.Run("BoundOnly", func(b *testing.B) {
		bound := bind.Float64UpDownCounter(base, userAlice, userID)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Add(ctx, 1.0)
			}
		})
	})

	b.Run("AddAddr", func(b *testing.B) {
		bound := bind.Float64UpDownCounter(base, userAlice, userID)
		extra := []metric.AddOption{metric.WithAttributes(adminTrue)}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Add(ctx, 1.0, extra...)
			}
		})
	})
}
