package bind

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Float64Histogram binds attrs to inst. All measurements made with the
// returned [metric.Float64Histogram] will include attrs.
//
// If inst is already bound to attributes, attrs will be merged into those
// attributes for the returned instrument.
func Float64Histogram(inst metric.Float64Histogram, attrs ...attribute.KeyValue) metric.Float64Histogram {
	if len(attrs) == 0 {
		return inst
	}

	if i, ok := inst.(float64Histogram); ok {
		// Flatten the instrument if already bound.
		inst = i.inst
		attrs = append(i.attrs, attrs...)
	}

	set := attribute.NewSet(attrs...)
	return float64Histogram{
		inst:  inst,
		attrs: attrs,
		set:   set,
		o:     []metric.RecordOption{metric.WithAttributeSet(set)},
	}
}

type float64Histogram struct {
	embedded.Float64Histogram

	inst  metric.Float64Histogram
	attrs []attribute.KeyValue
	set   attribute.Set
	o     []metric.RecordOption
}

// Unwrap returns the underlying [metric.Float64Histogram] and the bound
// attribute set.
func (i float64Histogram) Unwrap() (metric.Float64Histogram, attribute.Set) {
	return i.inst, i.set
}

// Record adds a value to the histogram. All measurements made will
// include the attributes bound to the instrument.
func (i float64Histogram) Record(ctx context.Context, value float64, opts ...metric.RecordOption) {
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
