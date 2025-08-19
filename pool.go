package bind

import (
	"sync"

	"go.opentelemetry.io/otel/metric"
)

var addOptPool = &sync.Pool{
	New: func() any {
		// This pool is used when adding additional attributes to a bound
		// instrument that already has attributes set. Meaning there should be
		// 2 options: WithAttributeSet + WithAttributes.
		s := make([]metric.AddOption, 0, 2)
		return &s
	},
}

var recordOptPool = &sync.Pool{
	New: func() any {
		// This pool is used for RecordOption slices for Record-based instruments.
		s := make([]metric.RecordOption, 0, 2)
		return &s
	},
}
