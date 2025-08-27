package bind_test

import (
	"testing"

	"github.com/MrAlias/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// mockMeter is a mock meter for testing.
type mockMeter struct {
	noop.Meter

	err error
}

func (m *mockMeter) Int64Counter(n string, o ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return errOr(mockInt64Counter{name: n, instOpts: o}, m.err)
}

func (m *mockMeter) Int64UpDownCounter(n string, o ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return errOr(mockInt64UpDownCounter{name: n, instOpts: o}, m.err)
}

func (m *mockMeter) Int64Histogram(n string, o ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return errOr(mockInt64Histogram{name: n, instOpts: o}, m.err)
}

func (m *mockMeter) Int64Gauge(n string, o ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return errOr(mockInt64Gauge{name: n, instOpts: o}, m.err)
}

func (m *mockMeter) Float64Counter(n string, o ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	return errOr(mockFloat64Counter{name: n, instOpts: o}, m.err)
}

func (m *mockMeter) Float64UpDownCounter(n string, o ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	return errOr(mockFloat64UpDownCounter{name: n, instOpts: o}, m.err)
}

func (m *mockMeter) Float64Histogram(n string, o ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return errOr(mockFloat64Histogram{name: n, instOpts: o}, m.err)
}

func (m *mockMeter) Float64Gauge(n string, o ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	return errOr(mockFloat64Gauge{name: n, instOpts: o}, m.err)
}

func errOr[T any](t T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Error-returning meter for testing error cases
func TestMeterEmptyAttrs(t *testing.T) {
	mock := &mockMeter{}
	got := bind.Meter(mock)
	assert.Same(t, mock, got, "meter with no attributes should return the same meter")
}

func TestMeterWithAttrs(t *testing.T) {
	mock := &mockMeter{}
	attrs := []attribute.KeyValue{userAlice, userID}

	bound := bind.Meter(mock, attrs...)
	require.NotNil(t, bound, "bound meter should not be nil")
	assert.NotSame(t, mock, bound, "bound meter should be different from original")

	in64Counter, err := bound.Int64Counter("test_counter")
	require.NoError(t, err)
	_, got := bind.Unwrap(in64Counter)
	assert.ElementsMatch(t, attrs, got.ToSlice(), "attributes should match bound attributes")
}

func TestMeterFlatten(t *testing.T) {
	mock := &mockMeter{}

	bound0 := bind.Meter(mock, userAlice)
	require.NotNil(t, bound0, "first bound meter should not be nil")

	bound1 := bind.Meter(bound0, userID)
	require.NotNil(t, bound1, "second bound meter should not be nil")

	got, _ := bind.Unwrap(bound1)
	assert.Same(t, mock, got, "underlying meter should be the original mock meter")
}

func TestMeterNoSideEffects(t *testing.T) {
	mock := &mockMeter{}

	a, cpA := clone(attribute.Int("C", 3), attribute.Int("B", 2))
	b, cpB := clone(attribute.Int("D", 4), attribute.Int("A", 1))

	m1 := bind.Meter(mock, a...)
	assert.Equal(t, cpA, a, "original attributes should not be modified")

	_ = bind.Meter(m1, b...)
	assert.Equal(t, cpA, a, "original attributes should not be modified after second bind")
	assert.Equal(t, cpB, b, "second set of attributes should not be modified")
}

func testMeterInst[T any](instFn func(m metric.Meter) (T, error)) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("CreateInstrument", func(t *testing.T) {
			attrs := []attribute.KeyValue{userAlice, userID}
			meter := bind.Meter(&mockMeter{}, attrs...)

			inst, err := instFn(meter)
			require.NoError(t, err)
			require.NotNil(t, inst, "instrument should not be nil")

			_, set := bind.Unwrap(inst)
			assert.ElementsMatch(t, attrs, set.ToSlice(), "unwrapped attributes should match bound attributes")
		})

		t.Run("CreateInstrumentError", func(t *testing.T) {
			meter := &mockMeter{err: assert.AnError}
			inst, err := instFn(meter)
			assert.ErrorIs(t, err, assert.AnError)
			assert.Nil(t, inst)
		})
	}
}

// Test all instrument creation methods
func TestMeterInstruments(t *testing.T) {
	const name = "test_instrument"

	t.Run("Int64Counter", testMeterInst(func(m metric.Meter) (metric.Int64Counter, error) {
		return m.Int64Counter(name)
	}))
	t.Run("Int64UpDownCounter", testMeterInst(func(m metric.Meter) (metric.Int64UpDownCounter, error) {
		return m.Int64UpDownCounter(name)
	}))
	t.Run("Int64Histogram", testMeterInst(func(m metric.Meter) (metric.Int64Histogram, error) {
		return m.Int64Histogram(name)
	}))
	t.Run("Int64Gauge", testMeterInst(func(m metric.Meter) (metric.Int64Gauge, error) {
		return m.Int64Gauge(name)
	}))

	t.Run("Float64Counter", testMeterInst(func(m metric.Meter) (metric.Float64Counter, error) {
		return m.Float64Counter(name)
	}))
	t.Run("Float64UpDownCounter", testMeterInst(func(m metric.Meter) (metric.Float64UpDownCounter, error) {
		return m.Float64UpDownCounter(name)
	}))
	t.Run("Float64Histogram", testMeterInst(func(m metric.Meter) (metric.Float64Histogram, error) {
		return m.Float64Histogram(name)
	}))
	t.Run("Float64Gauge", testMeterInst(func(m metric.Meter) (metric.Float64Gauge, error) {
		return m.Float64Gauge(name)
	}))
}

func BenchmarkMeter(b *testing.B) {
	run := func(m metric.Meter) func(*testing.B) {
		return func(b *testing.B) {
			var (
				iCntr   metric.Int64Counter
				iUDCntr metric.Int64UpDownCounter
				iHist   metric.Int64Histogram
				iGauge  metric.Int64Gauge

				iCntrErr, iUDCntrErr, iHistErr, iGaugeErr error

				fCntr   metric.Float64Counter
				fUDCntr metric.Float64UpDownCounter
				fHist   metric.Float64Histogram
				fGauge  metric.Float64Gauge

				fCntrErr, fUDCntrErr, fHistErr, fGaugeErr error
			)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				iCntr, iCntrErr = m.Int64Counter("int64 counter")
				iUDCntr, iUDCntrErr = m.Int64UpDownCounter("int64 updown")
				iHist, iHistErr = m.Int64Histogram("int64 histogram")
				iGauge, iGaugeErr = m.Int64Gauge("int64 gauge")

				fCntr, fCntrErr = m.Float64Counter("float64 counter")
				fUDCntr, fUDCntrErr = m.Float64UpDownCounter("float64 updown")
				fHist, fHistErr = m.Float64Histogram("float64 histogram")
				fGauge, fGaugeErr = m.Float64Gauge("float64 gauge")
			}

			_, _ = iCntr, iCntrErr
			_, _ = iUDCntr, iUDCntrErr
			_, _ = iHist, iHistErr
			_, _ = iGauge, iGaugeErr
			_, _ = fCntr, fCntrErr
			_, _ = fUDCntr, fUDCntrErr
			_, _ = fHist, fHistErr
			_, _ = fGauge, fGaugeErr
		}
	}

	base := noop.NewMeterProvider().Meter("test")
	b.Run("NoBinding", run(base))

	bound := bind.Meter(base, userAlice, userID)
	b.Run("WithBinding", run(bound))

	bound = bind.Meter(base, userAlice)
	bound = bind.Meter(bound, userID)
	b.Run("DoubleBinding", run(bound))
}
