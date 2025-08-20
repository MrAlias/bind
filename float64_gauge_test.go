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

type mockFloat64Gauge struct {
	embedded.Float64Gauge

	val  *float64
	opts []metric.RecordOption
}

func (m *mockFloat64Gauge) Record(ctx context.Context, val float64, opts ...metric.RecordOption) {
	m.val = &val
	m.opts = opts
}

func (m *mockFloat64Gauge) Instrument() metric.Float64Gauge {
	return m
}

func (m *mockFloat64Gauge) Recorded() (*float64, []attribute.KeyValue) {
	set := metric.NewRecordConfig(m.opts).Attributes()
	return m.val, set.ToSlice()
}

func TestFloat64Gauge(t *testing.T) {
	tests := []TestCase[float64]{
		{"BoundOnly", 3.14, nil},
		{"AddAttr", 2.71, []attribute.KeyValue{adminTrue}},
	}

	for _, test := range tests {
		t.Run(test.Name, Run(
			&mockFloat64Gauge{},
			bind.Float64Gauge,
			measFloat64Gauge,
			test,
		))
	}
}

func measFloat64Gauge(i metric.Float64Gauge, ctx context.Context, val float64, attr []attribute.KeyValue) {
	var opts []metric.RecordOption
	if attr != nil {
		opts = append(opts, metric.WithAttributes(attr...))
	}
	i.Record(ctx, val, opts...)
}

func TestFloat64GaugeUnwrap(t *testing.T) {
	mock := &mockFloat64Gauge{}
	bound := bind.Float64Gauge(mock, userAlice)
	require.NotNil(t, bound, "bound should not be nil")
	bound = bind.Float64Gauge(bound, userID)
	require.NotNil(t, bound, "bound should not be nil")

	val, set := bind.Unwrap(bound)

	assert.Same(t, mock, val, "unwrapped value should match mock")

	want := []attribute.KeyValue{userAlice, userID}
	assert.ElementsMatch(t, want, set.ToSlice(), "unwrapped attributes")
}

func BenchmarkFloat64GaugeRecord(b *testing.B) {
	ctx := context.Background()
	base := noop.Float64Gauge{}

	b.Run("BoundOnly", func(b *testing.B) {
		bound := bind.Float64Gauge(base, userAlice, userID)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Record(ctx, 1.0)
			}
		})
	})

	b.Run("AddAddr", func(b *testing.B) {
		bound := bind.Float64Gauge(base, userAlice, userID)
		extra := []metric.RecordOption{metric.WithAttributes(adminTrue)}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Record(ctx, 1.0, extra...)
			}
		})
	})
}
