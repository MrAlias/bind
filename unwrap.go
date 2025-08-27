package bind

import "go.opentelemetry.io/otel/attribute"

type unwrapper[T any] interface {
	Unwrap() (T, attribute.Set)
}

// Unwrap unwraps any bound instrument returning the unwrapped instrument and
// any attributes that it was bound to.
func Unwrap[T any](inst T) (T, attribute.Set) {
	if u, ok := any(inst).(unwrapper[T]); ok {
		return u.Unwrap()
	}
	return inst, *attribute.EmptySet()
}
