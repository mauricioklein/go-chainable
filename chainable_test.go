package chainable

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	sum2 := func(x int) int { return x + 2 }
	mul2 := func(x int) int { return x * 2 }

	// (4 + 2) * 2 = 12
	ret, err := New().
		From(4).
		Chain(sum2, mul2).
		Unwrap()

	assert.Equal(t, []interface{}{12}, ret)
	assert.NoError(t, err)
}

func TestChainDummy(t *testing.T) {
	genericError := errors.New("a generic error")

	f1 := func() error { return genericError }
	f2 := func(e error) int { return 0 }

	_, err := New().ChainDummy(f1, f2).Unwrap()

	assert.NoError(t, err)
}

func TestNotAFunction(t *testing.T) {
	_, err := New().Chain("not a function").Unwrap()
	assert.EqualError(t, err, errNotAFunction.Error())
}

func TestArgumentMismatch(t *testing.T) {
	f1 := func() (int, error) { return 0, nil }
	f2 := func(a, b int) (int, error) { return 0, nil }

	_, err := New().Chain(f1, f2).Unwrap()
	assert.EqualError(t, err, errArgumentMismatch.Error())
}

func TestReset(t *testing.T) {
	returnOne := func() int { return 1 }

	res, err := New().Chain(returnOne).Reset().Unwrap()

	assert.Equal(t, []interface{}{}, res)
	assert.NoError(t, err)
}
