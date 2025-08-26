package bind

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Float64Gauge binds attrs to inst. All measurements made with the
// returned [metric.Float64Gauge] will include attrs.
//
// If inst is already bound to attributes, attrs will be merged into those
// attributes for the returned instrument.
func Float64Gauge(inst metric.Float64Gauge, attrs ...attribute.KeyValue) metric.Float64Gauge {
	if len(attrs) == 0 {
		return inst
	}

	// NewSet sorts passed attributes. Copy to avoid side effect.
	var cp []attribute.KeyValue

	if i, ok := inst.(float64Gauge); ok {
		// Flatten the instrument if already bound.
		inst = i.inst

		cp = make([]attribute.KeyValue, 0, len(i.attrs)+len(attrs))
		cp = append(cp, i.attrs...)
		cp = append(cp, attrs...)
	} else {
		cp = make([]attribute.KeyValue, len(attrs))
		copy(cp, attrs)
	}

	set := attribute.NewSet(cp...)
	return float64Gauge{
		inst:  inst,
		attrs: cp,
		set:   set,
		o:     []metric.RecordOption{metric.WithAttributeSet(set)},
	}
}

type float64Gauge struct {
	embedded.Float64Gauge

	inst  metric.Float64Gauge
	attrs []attribute.KeyValue
	set   attribute.Set
	o     []metric.RecordOption
}

// Unwrap returns the underlying [metric.Float64Gauge] and the bound
// attribute set.
func (i float64Gauge) Unwrap() (metric.Float64Gauge, attribute.Set) {
	return i.inst, i.set
}

// Record records the instantaneous value. All measurements made will
// include the attributes bound to the instrument.
func (i float64Gauge) Record(ctx context.Context, value float64, opts ...metric.RecordOption) {
	if len(opts) == 0 {
		i.inst.Record(ctx, value, i.o...)
		return
	}

	o := recordOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recordOptPool.Put(o)
	}()

	*o = append(*o, i.o...)
	*o = append(*o, opts...)
	i.inst.Record(ctx, value, *o...)
}
