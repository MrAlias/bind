package bind

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Float64UpDownCounter binds attrs to inst. All measurements made with the
// returned [metric.Float64UpDownCounter] will include attrs.
//
// If inst is already bound to attributes, attrs will be merged into those
// attributes for the returned instrument.
func Float64UpDownCounter(inst metric.Float64UpDownCounter, attrs ...attribute.KeyValue) metric.Float64UpDownCounter {
	if i, ok := inst.(float64UpDownCounter); ok {
		// Flatten the instrument if already bound.
		inst = i.inst
		attrs = append(i.attrs, attrs...)
	}

	set := attribute.NewSet(attrs...)
	return float64UpDownCounter{
		inst:  inst,
		attrs: attrs,
		set:   set,
		o:     []metric.AddOption{metric.WithAttributeSet(set)},
	}
}

type float64UpDownCounter struct {
	embedded.Float64UpDownCounter

	inst  metric.Float64UpDownCounter
	attrs []attribute.KeyValue
	set   attribute.Set
	o     []metric.AddOption
}

// Unwrap returns the underlying [metric.Float64UpDownCounter] and the bound
// attribute set.
func (i float64UpDownCounter) Unwrap() (metric.Float64UpDownCounter, attribute.Set) {
	return i.inst, i.set
}

// Add records a change to the counter. All measurements made will
// include the attributes bound to the instrument.
func (i float64UpDownCounter) Add(ctx context.Context, incr float64, opts ...metric.AddOption) {
	if len(opts) == 0 {
		i.inst.Add(ctx, incr, i.o...)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, i.o...)
	*o = append(*o, opts...)
	i.inst.Add(ctx, incr, *o...)
}
