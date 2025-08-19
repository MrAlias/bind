package bind

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Int64Gauge binds attrs to inst. All measurements made with the
// returned [metric.Int64Gauge] will include attrs.
//
// If inst is already bound to attributes, attrs will be merged into those
// attributes for the returned instrument.
func Int64Gauge(inst metric.Int64Gauge, attrs ...attribute.KeyValue) metric.Int64Gauge {
	if i, ok := inst.(int64Gauge); ok {
		// Flatten the instrument if already bound.
		inst = i.inst
		attrs = append(i.attrs, attrs...)
	}

	set := attribute.NewSet(attrs...)
	return int64Gauge{
		inst:  inst,
		attrs: attrs,
		set:   set,
		o:     []metric.RecordOption{metric.WithAttributeSet(set)},
	}
}

type int64Gauge struct {
	embedded.Int64Gauge

	inst  metric.Int64Gauge
	attrs []attribute.KeyValue
	set   attribute.Set
	o     []metric.RecordOption
}

// Unwrap returns the underlying [metric.Int64Gauge] and the bound
// attribute set.
func (i int64Gauge) Unwrap() (metric.Int64Gauge, attribute.Set) {
	return i.inst, i.set
}

// Record records the instantaneous value. All measurements made will
// include the attributes bound to the instrument.
func (i int64Gauge) Record(ctx context.Context, value int64, opts ...metric.RecordOption) {
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
