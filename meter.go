package bind

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Meter binds attrs to m. All synchronous instruments created with the
// returned [metric.Meter] will bind attrs to the equivalent synchronous
// instrument created by m.
func Meter(m metric.Meter, attrs ...attribute.KeyValue) metric.Meter {
	if len(attrs) == 0 {
		return m
	}

	// NewSet sorts passed attributes. Copy to avoid side effect.
	var cp []attribute.KeyValue

	if i, ok := m.(*meter); ok {
		// Flatten the meter if already bound.
		m = i.Meter
		cp = make([]attribute.KeyValue, 0, len(i.attrs)+len(attrs))
		cp = append(cp, i.attrs...)
		cp = append(cp, attrs...)
	} else {
		cp = make([]attribute.KeyValue, len(attrs))
		copy(cp, attrs)
	}

	set := attribute.NewSet(cp...)
	o := metric.WithAttributeSet(set)
	return &meter{
		Meter:  m,
		attrs:  cp,
		set:    set,
		addOpt: []metric.AddOption{o},
		recOpt: []metric.RecordOption{o},
	}
}

type meter struct {
	metric.Meter

	attrs  []attribute.KeyValue
	set    attribute.Set
	addOpt []metric.AddOption
	recOpt []metric.RecordOption
}

var (
	_ metric.Meter            = (*meter)(nil)
	_ unwrapper[metric.Meter] = (*meter)(nil)
)

func (m *meter) Unwrap() (metric.Meter, attribute.Set) {
	return m.Meter, m.set
}

func (m *meter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	inst, err := m.Meter.Int64Counter(name, options...)
	if inst != nil {
		inst = int64Counter{
			inst:  inst,
			attrs: m.attrs,
			set:   m.set,
			o:     m.addOpt,
		}
	}
	return inst, err
}

func (m *meter) Int64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	inst, err := m.Meter.Int64UpDownCounter(name, options...)
	if inst != nil {
		inst = int64UpDownCounter{
			inst:  inst,
			attrs: m.attrs,
			set:   m.set,
			o:     m.addOpt,
		}
	}
	return inst, err
}

func (m *meter) Int64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	inst, err := m.Meter.Int64Histogram(name, options...)
	if inst != nil {
		inst = int64Histogram{
			inst:  inst,
			attrs: m.attrs,
			set:   m.set,
			o:     m.recOpt,
		}
	}
	return inst, err
}

func (m *meter) Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	inst, err := m.Meter.Int64Gauge(name, options...)
	if inst != nil {
		inst = int64Gauge{
			inst:  inst,
			attrs: m.attrs,
			set:   m.set,
			o:     m.recOpt,
		}
	}
	return inst, err
}

func (m *meter) Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	inst, err := m.Meter.Float64Counter(name, options...)
	if inst != nil {
		inst = float64Counter{
			inst:  inst,
			attrs: m.attrs,
			set:   m.set,
			o:     m.addOpt,
		}
	}
	return inst, err
}

func (m *meter) Float64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	inst, err := m.Meter.Float64UpDownCounter(name, options...)
	if inst != nil {
		inst = float64UpDownCounter{
			inst:  inst,
			attrs: m.attrs,
			set:   m.set,
			o:     m.addOpt,
		}
	}
	return inst, err
}

func (m *meter) Float64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	inst, err := m.Meter.Float64Histogram(name, options...)
	if inst != nil {
		inst = float64Histogram{
			inst:  inst,
			attrs: m.attrs,
			set:   m.set,
			o:     m.recOpt,
		}
	}
	return inst, err
}

func (m *meter) Float64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	inst, err := m.Meter.Float64Gauge(name, options...)
	if inst != nil {
		inst = float64Gauge{
			inst:  inst,
			attrs: m.attrs,
			set:   m.set,
			o:     m.recOpt,
		}
	}
	return inst, err
}
