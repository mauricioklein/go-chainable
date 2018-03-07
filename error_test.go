package chainable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotAFunctionError(t *testing.T) {
	err := NewNotAFunctionError(1)
	assert.EqualError(t, err, "(Error on Link: 1) Element isn't a function")
}

func TestArgumentMismatchError(t *testing.T) {
	err := NewArgumentMismatchError(1, 2, 3)
	assert.EqualError(t, err, "(Error on Link: 1) 2 arg(s) provided, but function arity is 3")
}
