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

type mockInt64Gauge struct {
	embedded.Int64Gauge

	name     string
	instOpts []metric.Int64GaugeOption

	val     *int64
	recOpts []metric.RecordOption
}

func (m *mockInt64Gauge) Record(ctx context.Context, val int64, opts ...metric.RecordOption) {
	m.val = &val
	m.recOpts = opts
}

func (m *mockInt64Gauge) Instrument() metric.Int64Gauge {
	return m
}

func (m *mockInt64Gauge) Recorded() (*int64, []attribute.KeyValue) {
	set := metric.NewRecordConfig(m.recOpts).Attributes()
	return m.val, set.ToSlice()
}

func TestInt64GaugeEmptyAttrs(t *testing.T) {
	mock := &mockInt64Gauge{}
	got := bind.Int64Gauge(mock)
	assert.Same(t, mock, got, "bound should be the same as the input")
}

func TestInt64Gauge(t *testing.T) {
	tests := []TestCase[int64]{
		{"BoundOnly", 11, nil},
		{"AddAttr", 22, []attribute.KeyValue{adminTrue}},
	}

	for _, test := range tests {
		t.Run(test.Name, Run(
			&mockInt64Gauge{},
			bind.Int64Gauge,
			measInt64Gauge,
			test,
		))
	}

	t.Run(
		"NoSideEffets",
		testNoSideEffects(bind.Int64Gauge, &mockInt64Gauge{}),
	)
}

func measInt64Gauge(i metric.Int64Gauge, ctx context.Context, val int64, attr []attribute.KeyValue) {
	var opts []metric.RecordOption
	if attr != nil {
		opts = append(opts, metric.WithAttributes(attr...))
	}
	i.Record(ctx, val, opts...)
}

func TestInt64GaugeUnwrap(t *testing.T) {
	mock := &mockInt64Gauge{}
	bound := bind.Int64Gauge(mock, userAlice)
	require.NotNil(t, bound, "bound should not be nil")
	bound = bind.Int64Gauge(bound, userID)
	require.NotNil(t, bound, "bound should not be nil")

	val, set := bind.Unwrap(bound)

	assert.Same(t, mock, val, "unwrapped value should match mock")

	want := []attribute.KeyValue{userAlice, userID}
	assert.ElementsMatch(t, want, set.ToSlice(), "unwrapped attributes")
}

func BenchmarkInt64GaugeRecord(b *testing.B) {
	ctx := context.Background()
	base := noop.Int64Gauge{}

	b.Run("BoundOnly", func(b *testing.B) {
		bound := bind.Int64Gauge(base, userAlice, userID)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Record(ctx, 1.0)
			}
		})
	})

	b.Run("AddAddr", func(b *testing.B) {
		bound := bind.Int64Gauge(base, userAlice, userID)
		extra := []metric.RecordOption{metric.WithAttributes(adminTrue)}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bound.Record(ctx, 1.0, extra...)
			}
		})
	})
}
