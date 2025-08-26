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

type mockInt64Histogram struct {
	embedded.Int64Histogram

	val  *int64
	opts []metric.RecordOption
}

func (m *mockInt64Histogram) Record(ctx context.Context, val int64, opts ...metric.RecordOption) {
	m.val = &val
	m.opts = opts
}

func (m *mockInt64Histogram) Instrument() metric.Int64Histogram {
	return m
}

func (m *mockInt64Histogram) Recorded() (*int64, []attribute.KeyValue) {
	set := metric.NewRecordConfig(m.opts).Attributes()
	return m.val, set.ToSlice()
}

func TestInt64HistogramEmptyAttrs(t *testing.T) {
	mock := &mockInt64Histogram{}
	got := bind.Int64Histogram(mock)
	assert.Same(t, mock, got, "bound should be the same as the input")
}

func TestInt64Histogram(t *testing.T) {
	tests := []TestCase[int64]{
		{"BoundOnly", 11, nil},
		{"AddAttr", 22, []attribute.KeyValue{adminTrue}},
	}

	for _, test := range tests {
		t.Run(test.Name, Run(
			&mockInt64Histogram{},
			bind.Int64Histogram,
			measInt64Histogram,
			test,
		))
	}
}

func measInt64Histogram(i metric.Int64Histogram, ctx context.Context, val int64, attr []attribute.KeyValue) {
	var opts []metric.RecordOption
	if attr != nil {
		opts = append(opts, metric.WithAttributes(attr...))
	}
	i.Record(ctx, val, opts...)
}

func TestInt64HistogramUnwrap(t *testing.T) {
	mock := &mockInt64Histogram{}
	bound := bind.Int64Histogram(mock, userAlice)
	require.NotNil(t, bound, "bound should not be nil")
	bound = bind.Int64Histogram(bound, userID)
	require.NotNil(t, bound, "bound should not be nil")

	val, set := bind.Unwrap(bound)

	assert.Same(t, mock, val, "unwrapped value should match mock")

	want := []attribute.KeyValue{userAlice, userID}
	assert.ElementsMatch(t, want, set.ToSlice(), "unwrapped attributes")
}

func BenchmarkInt64HistogramRecord(b *testing.B) {
	ctx := context.Background()
	base := noop.Int64Histogram{}

	b.Run("BoundOnly", func(b *testing.B) {
		bound := bind.Int64Histogram(base, userAlice, userID)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Record(ctx, 1.0)
			}
		})
	})

	b.Run("AddAddr", func(b *testing.B) {
		bound := bind.Int64Histogram(base, userAlice, userID)
		extra := []metric.RecordOption{metric.WithAttributes(adminTrue)}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Record(ctx, 1.0, extra...)
			}
		})
	})
}
