package bind

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Int64UpDownCounter binds attrs to inst. All measurements made with the
// returned [metric.Int64UpDownCounter] will include attrs.
//
// If inst is already bound to attributes, attrs will be merged into those
// attributes for the returned instrument.
func Int64UpDownCounter(inst metric.Int64UpDownCounter, attrs ...attribute.KeyValue) metric.Int64UpDownCounter {
	if len(attrs) == 0 {
		return inst
	}

	if i, ok := inst.(int64UpDownCounter); ok {
		// Flatten the instrument if already bound.
		inst = i.inst
		attrs = append(i.attrs, attrs...)
	}

	set := attribute.NewSet(attrs...)
	return int64UpDownCounter{
		inst:  inst,
		attrs: attrs,
		set:   set,
		o:     []metric.AddOption{metric.WithAttributeSet(set)},
	}
}

type int64UpDownCounter struct {
	embedded.Int64UpDownCounter

	inst  metric.Int64UpDownCounter
	attrs []attribute.KeyValue
	set   attribute.Set
	o     []metric.AddOption
}

// Unwrap returns the underlying [metric.Int64UpDownCounter] and the bound
// attribute set.
func (i int64UpDownCounter) Unwrap() (metric.Int64UpDownCounter, attribute.Set) {
	return i.inst, i.set
}

// Add increments or decrements the counter by incr. All measurements made will
// include the attributes bound to the instrument.
func (i int64UpDownCounter) Add(ctx context.Context, incr int64, opts ...metric.AddOption) {
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
