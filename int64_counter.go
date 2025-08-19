package bind

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Int64Counter binds attrs to inst. All measurements made with the
// returned [metric.Int64Counter] will include attrs.
//
// If inst is already bound to attributes, attrs will be merged into those
// attributes for the returned instrument.
func Int64Counter(inst metric.Int64Counter, attrs ...attribute.KeyValue) metric.Int64Counter {
	if i, ok := inst.(int64Counter); ok {
		// Flatten the instrument if already bound.
		inst = i.inst
		attrs = append(i.attrs, attrs...)
	}

	set := attribute.NewSet(attrs...)
	return int64Counter{
		inst:  inst,
		attrs: attrs,
		set:   set,
		o:     []metric.AddOption{metric.WithAttributeSet(set)},
	}
}

type int64Counter struct {
	embedded.Int64Counter

	inst  metric.Int64Counter
	attrs []attribute.KeyValue
	set   attribute.Set
	o     []metric.AddOption
}

// Unwrap returns the underlying [metric.Int64Counter] and the bound
// attribute set.
func (i int64Counter) Unwrap() (metric.Int64Counter, attribute.Set) {
	return i.inst, i.set
}

// Add increments the counter by incr. All measurements made will
// include the attributes bound to the instrument.
func (i int64Counter) Add(ctx context.Context, incr int64, opts ...metric.AddOption) {
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
