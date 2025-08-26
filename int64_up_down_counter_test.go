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

type mockInt64UpDownCounter struct {
	embedded.Int64UpDownCounter

	incr *int64
	opts []metric.AddOption
}

func (m *mockInt64UpDownCounter) Add(ctx context.Context, incr int64, opts ...metric.AddOption) {
	m.incr = &incr
	m.opts = opts
}

func (m *mockInt64UpDownCounter) Instrument() metric.Int64UpDownCounter {
	return m
}

func (m *mockInt64UpDownCounter) Recorded() (*int64, []attribute.KeyValue) {
	set := metric.NewAddConfig(m.opts).Attributes()
	return m.incr, set.ToSlice()
}

func TestInt64UpDownCounterEmptyAttrs(t *testing.T) {
	mock := &mockInt64UpDownCounter{}
	got := bind.Int64UpDownCounter(mock)
	assert.Same(t, mock, got, "bound should be the same as the input")
}

func TestInt64UpDownCounter(t *testing.T) {
	tests := []TestCase[int64]{
		{"BoundOnly", 100, nil},
		{"AddAttr", 200, []attribute.KeyValue{adminTrue}},
	}

	for _, test := range tests {
		t.Run(test.Name, Run(
			&mockInt64UpDownCounter{},
			bind.Int64UpDownCounter,
			measInt64UpDownCounter,
			test,
		))
	}
}

func measInt64UpDownCounter(i metric.Int64UpDownCounter, ctx context.Context, incr int64, attr []attribute.KeyValue) {
	var opts []metric.AddOption
	if attr != nil {
		opts = append(opts, metric.WithAttributes(attr...))
	}
	i.Add(ctx, incr, opts...)
}

func TestInt64UpDownCounterUnwrap(t *testing.T) {
	mock := &mockInt64UpDownCounter{}
	bound := bind.Int64UpDownCounter(mock, userAlice)
	require.NotNil(t, bound, "bound should not be nil")
	bound = bind.Int64UpDownCounter(bound, userID)
	require.NotNil(t, bound, "bound should not be nil")

	val, set := bind.Unwrap(bound)

	assert.Same(t, mock, val, "unwrapped value should match mock")

	want := []attribute.KeyValue{userAlice, userID}
	assert.ElementsMatch(t, want, set.ToSlice(), "unwrapped attributes")
}

func BenchmarkInt64UpDownCounterAdd(b *testing.B) {
	ctx := context.Background()
	base := noop.Int64UpDownCounter{}

	b.Run("BoundOnly", func(b *testing.B) {
		bound := bind.Int64UpDownCounter(base, userAlice, userID)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Add(ctx, 1.0)
			}
		})
	})

	b.Run("AddAddr", func(b *testing.B) {
		bound := bind.Int64UpDownCounter(base, userAlice, userID)
		extra := []metric.AddOption{metric.WithAttributes(adminTrue)}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Add(ctx, 1.0, extra...)
			}
		})
	})
}
