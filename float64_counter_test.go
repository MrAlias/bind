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

type mockFloat64Counter struct {
	embedded.Float64Counter

	name     string
	instOpts []metric.Float64CounterOption

	incr    *float64
	addOpts []metric.AddOption
}

func (m *mockFloat64Counter) Add(ctx context.Context, incr float64, opts ...metric.AddOption) {
	m.incr = &incr
	m.addOpts = opts
}

func (m *mockFloat64Counter) Instrument() metric.Float64Counter {
	return m
}

func (m *mockFloat64Counter) Recorded() (*float64, []attribute.KeyValue) {
	set := metric.NewAddConfig(m.addOpts).Attributes()
	return m.incr, set.ToSlice()
}

func TestFloat64CounterEmptyAttrs(t *testing.T) {
	mock := &mockFloat64Counter{}
	got := bind.Float64Counter(mock)
	assert.Same(t, mock, got, "bound should be the same as the input")
}

func TestFloat64Counter(t *testing.T) {
	tests := []TestCase[float64]{
		{"BoundOnly", 42.0, nil},
		{"AddAttr", 7.0, []attribute.KeyValue{adminTrue}},
	}

	for _, test := range tests {
		t.Run(test.Name, Run(
			&mockFloat64Counter{},
			bind.Float64Counter,
			measFloat64Counter,
			test,
		))
	}

	t.Run(
		"NoSideEffets",
		testNoSideEffects(bind.Float64Counter, &mockFloat64Counter{}),
	)
}

func measFloat64Counter(i metric.Float64Counter, ctx context.Context, incr float64, attr []attribute.KeyValue) {
	var opts []metric.AddOption
	if attr != nil {
		opts = append(opts, metric.WithAttributes(attr...))
	}
	i.Add(ctx, incr, opts...)
}

func TestFloat64CounterUnwrap(t *testing.T) {
	mock := &mockFloat64Counter{}
	bound := bind.Float64Counter(mock, userAlice)
	require.NotNil(t, bound, "bound should not be nil")
	bound = bind.Float64Counter(bound, userID)
	require.NotNil(t, bound, "bound should not be nil")

	val, set := bind.Unwrap(bound)

	assert.Same(t, mock, val, "unwrapped value should match mock")

	want := []attribute.KeyValue{userAlice, userID}
	assert.ElementsMatch(t, want, set.ToSlice(), "unwrapped attributes")
}

func BenchmarkFloat64CounterAdd(b *testing.B) {
	ctx := context.Background()
	base := noop.Float64Counter{}

	b.Run("BoundOnly", func(b *testing.B) {
		bound := bind.Float64Counter(base, userAlice, userID)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Add(ctx, 1.0)
			}
		})
	})

	b.Run("AddAddr", func(b *testing.B) {
		bound := bind.Float64Counter(base, userAlice, userID)
		extra := []metric.AddOption{metric.WithAttributes(adminTrue)}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Add(ctx, 1.0, extra...)
			}
		})
	})
}
