package bind_test

import (
	"testing"

	"github.com/MrAlias/bind"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestUnwrapNonBound(t *testing.T) {
	mock := &mockInt64Counter{}

	unwrapped, attrs := bind.Unwrap(mock)
	assert.Same(t, mock, unwrapped, "unwrapped non-bound instrument should be the same")
	assert.Equal(t, *attribute.EmptySet(), attrs, "non-bound instrument should have empty attribute set")
}
