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
		From([]interface{}{4}).
		Chain(sum2).
		Chain(mul2).
		Unwrap()

	assert.Equal(t, []interface{}{12}, ret)
	assert.Nil(t, err)
}

func TestChain_NotAFunction(t *testing.T) {
	_, err := New().Chain("not a function").Unwrap()
	assert.EqualError(t, err, errNotAFunction.Error())
}

func TestChain_ArgumentMismatch(t *testing.T) {
	f1 := func() (int, error) { return 0, nil }
	f2 := func(a, b int) (int, error) { return 0, nil }

	_, err := New().Chain(f1).Chain(f2).Unwrap()
	assert.EqualError(t, err, errArgumentMismatch.Error())
}

func TestChain_RightmostValueNotAnError(t *testing.T) {
	f := func() int { return 0 }

	_, err := New().Chain(f).Unwrap()
	assert.EqualError(t, err, errRightmostValueNotAnError.Error())
}
