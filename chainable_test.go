package chainable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain_Success(t *testing.T) {
	sum2 := func(x int) (int, error) { return x + 2, nil }
	mul2 := func(x int) (int, error) { return x * 2, nil }

	// (4 + 2) * 2 = 12
	ret, err := New().
		From(4).
		Link(sum2, mul2).
		Unwrap()

	assert.Equal(t, []interface{}{12}, ret)
	assert.NoError(t, err)
}

func TestChain_NotAFunction(t *testing.T) {
	_, err := New().Link("not a function").Unwrap()
	assert.EqualError(t, err, errNotAFunction.Error())
}

func TestChain_ArgumentMismatch(t *testing.T) {
	f1 := func() (int, error) { return 0, nil }
	f2 := func(a, b int) (int, error) { return 0, nil }

	_, err := New().Link(f1, f2).Unwrap()
	assert.EqualError(t, err, errArgumentMismatch.Error())
}

func TestChain_AcceptFunctionNotReturningError(t *testing.T) {
	f1 := func() int { return 0 }

	_, err := New().Link(f1).Unwrap()
	assert.NoError(t, err)
}
