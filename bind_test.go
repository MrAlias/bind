package bind_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

var (
	userAlice = attribute.String("user", "alice")
	userID    = attribute.Int("id", 12345)
	adminTrue = attribute.Bool("admin", true)
)

type TestCase[N any] struct {
	Name     string
	Measure  N
	AddAttrs []attribute.KeyValue
}

type Mock[T any, N any] interface {
	Instrument() T
	Recorded() (*N, []attribute.KeyValue)
}

type Binder[T any] func(T, ...attribute.KeyValue) T

type Measure[T any, N any] func(T, context.Context, N, []attribute.KeyValue)

func Run[T any, N any](mock Mock[T, N], b Binder[T], m Measure[T, N], test TestCase[N]) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		bound := b(mock.Instrument(), userAlice)
		require.NotNil(t, bound, "bound instrument nil")

		bound = b(bound, userID)
		require.NotNil(t, bound, "second bound instrument nil")

		m(bound, context.Background(), test.Measure, test.AddAttrs)

		valPtr, attrs := mock.Recorded()

		if assert.NotNil(t, valPtr, "measure value not set") {
			assert.Equal(t, test.Measure, *valPtr, "measure value")
		}

		want := []attribute.KeyValue{userAlice, userID}
		want = append(want, test.AddAttrs...)
		assert.ElementsMatch(t, want, attrs, "attributes in config")
	}
}
