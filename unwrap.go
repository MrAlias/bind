package bind

import "go.opentelemetry.io/otel/attribute"

// Unwrap unwraps any bound instrument returning the unwrapped instrument and
// any attributes that it was bound to.
func Unwrap[T any](inst T) (T, attribute.Set) {
	type unwrapper interface {
		Unwrap() (T, attribute.Set)
	}

	if u, ok := any(inst).(unwrapper); ok {
		return u.Unwrap()
	}
	return inst, *attribute.EmptySet()
}
