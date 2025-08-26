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

type mockInt64Counter struct {
	embedded.Int64Counter

	incr *int64
	opts []metric.AddOption
}

func (m *mockInt64Counter) Add(ctx context.Context, incr int64, opts ...metric.AddOption) {
	m.incr = &incr
	m.opts = opts
}

func (m *mockInt64Counter) Instrument() metric.Int64Counter {
	return m
}

func (m *mockInt64Counter) Recorded() (*int64, []attribute.KeyValue) {
	set := metric.NewAddConfig(m.opts).Attributes()
	return m.incr, set.ToSlice()
}

func TestInt64CounterEmptyAttrs(t *testing.T) {
	mock := &mockInt64Counter{}
	got := bind.Int64Counter(mock)
	assert.Same(t, mock, got, "bound should be the same as the input")
}

func TestInt64Counter(t *testing.T) {
	tests := []TestCase[int64]{
		{"BoundOnly", 100, nil},
		{"AddAttr", 200, []attribute.KeyValue{adminTrue}},
	}

	for _, test := range tests {
		t.Run(test.Name, Run(
			&mockInt64Counter{},
			bind.Int64Counter,
			measInt64Counter,
			test,
		))
	}

	t.Run(
		"NoSideEffets",
		testNoSideEffects(bind.Int64Counter, &mockInt64Counter{}),
	)
}

func measInt64Counter(i metric.Int64Counter, ctx context.Context, incr int64, attr []attribute.KeyValue) {
	var opts []metric.AddOption
	if attr != nil {
		opts = append(opts, metric.WithAttributes(attr...))
	}
	i.Add(ctx, incr, opts...)
}

func TestInt64CounterUnwrap(t *testing.T) {
	mock := &mockInt64Counter{}
	bound := bind.Int64Counter(mock, userAlice)
	require.NotNil(t, bound, "bound should not be nil")
	bound = bind.Int64Counter(bound, userID)
	require.NotNil(t, bound, "bound should not be nil")

	val, set := bind.Unwrap(bound)

	assert.Same(t, mock, val, "unwrapped value should match mock")

	want := []attribute.KeyValue{userAlice, userID}
	assert.ElementsMatch(t, want, set.ToSlice(), "unwrapped attributes")
}

func BenchmarkInt64CounterAdd(b *testing.B) {
	ctx := context.Background()
	base := noop.Int64Counter{}

	b.Run("BoundOnly", func(b *testing.B) {
		bound := bind.Int64Counter(base, userAlice, userID)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Add(ctx, 1.0)
			}
		})
	})

	b.Run("AddAddr", func(b *testing.B) {
		bound := bind.Int64Counter(base, userAlice, userID)
		extra := []metric.AddOption{metric.WithAttributes(adminTrue)}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Add(ctx, 1.0, extra...)
			}
		})
	})
}
