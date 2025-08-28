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

type mockFloat64Histogram struct {
	embedded.Float64Histogram

	name     string
	instOpts []metric.Float64HistogramOption

	val     *float64
	recOpts []metric.RecordOption
}

func (m *mockFloat64Histogram) Record(_ context.Context, val float64, opts ...metric.RecordOption) {
	m.val = &val
	m.recOpts = opts
}

func (m *mockFloat64Histogram) Instrument() metric.Float64Histogram {
	return m
}

func (m *mockFloat64Histogram) Recorded() (*float64, []attribute.KeyValue) {
	set := metric.NewRecordConfig(m.recOpts).Attributes()
	return m.val, set.ToSlice()
}

func TestFloat64HistogramEmptyAttrs(t *testing.T) {
	mock := &mockFloat64Histogram{}
	got := bind.Float64Histogram(mock)
	assert.Same(t, mock, got, "bound should be the same as the input")
}

func TestFloat64Histogram(t *testing.T) {
	tests := []TestCase[float64]{
		{"BoundOnly", 1.23, nil},
		{"AddAttr", 4.56, []attribute.KeyValue{adminTrue}},
	}

	for _, test := range tests {
		t.Run(test.Name, Run(
			&mockFloat64Histogram{},
			bind.Float64Histogram,
			measFloat64Histogram,
			test,
		))
	}

	t.Run(
		"NoSideEffets",
		testNoSideEffects(bind.Float64Histogram, &mockFloat64Histogram{}),
	)
}

func measFloat64Histogram(i metric.Float64Histogram, ctx context.Context, val float64, attr []attribute.KeyValue) {
	var opts []metric.RecordOption
	if attr != nil {
		opts = append(opts, metric.WithAttributes(attr...))
	}
	i.Record(ctx, val, opts...)
}

func TestFloat64HistogramUnwrap(t *testing.T) {
	mock := &mockFloat64Histogram{}
	bound := bind.Float64Histogram(mock, userAlice)
	require.NotNil(t, bound, "bound should not be nil")
	bound = bind.Float64Histogram(bound, userID)
	require.NotNil(t, bound, "bound should not be nil")

	val, set := bind.Unwrap(bound)

	assert.Same(t, mock, val, "unwrapped value should match mock")

	want := []attribute.KeyValue{userAlice, userID}
	assert.ElementsMatch(t, want, set.ToSlice(), "unwrapped attributes")
}

func BenchmarkFloat64HistogramRecord(b *testing.B) {
	ctx := context.Background()
	base := noop.Float64Histogram{}

	b.Run("BoundOnly", func(b *testing.B) {
		bound := bind.Float64Histogram(base, userAlice, userID)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Record(ctx, 1.0)
			}
		})
	})

	b.Run("AddAddr", func(b *testing.B) {
		bound := bind.Float64Histogram(base, userAlice, userID)
		extra := []metric.RecordOption{metric.WithAttributes(adminTrue)}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Record(ctx, 1.0, extra...)
			}
		})
	})
}
