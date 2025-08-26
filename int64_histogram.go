package bind

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Int64Histogram binds attrs to inst. All measurements made with the
// returned [metric.Int64Histogram] will include attrs.
//
// If inst is already bound to attributes, attrs will be merged into those
// attributes for the returned instrument.
func Int64Histogram(inst metric.Int64Histogram, attrs ...attribute.KeyValue) metric.Int64Histogram {
	if len(attrs) == 0 {
		return inst
	}

	// NewSet sorts passed attributes. Copy to avoid side effect.
	var cp []attribute.KeyValue

	if i, ok := inst.(int64Histogram); ok {
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
	return int64Histogram{
		inst:  inst,
		attrs: cp,
		set:   set,
		o:     []metric.RecordOption{metric.WithAttributeSet(set)},
	}
}

type int64Histogram struct {
	embedded.Int64Histogram

	inst  metric.Int64Histogram
	attrs []attribute.KeyValue
	set   attribute.Set
	o     []metric.RecordOption
}

// Unwrap returns the underlying [metric.Int64Histogram] and the bound
// attribute set.
func (i int64Histogram) Unwrap() (metric.Int64Histogram, attribute.Set) {
	return i.inst, i.set
}

// Record adds a value to the histogram. All measurements made will
// include the attributes bound to the instrument.
func (i int64Histogram) Record(ctx context.Context, value int64, opts ...metric.RecordOption) {
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
