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
