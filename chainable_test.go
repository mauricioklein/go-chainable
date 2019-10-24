package chainable

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	testCases := []struct {
		desc        string
		from        []Argument
		funcs       []Function
		returnValue []Argument
		err         error
	}{
		{
			desc: "Regular chain",
			from: []Argument{4},
			funcs: []Function{
				func(x int) int { return x + 2 },
				func(x int) int { return x * 2 },
			},
			returnValue: []Argument{12},
			err:         nil,
		},
		{
			desc: "No initial value in chain",
			from: []Argument{},
			funcs: []Function{
				func() int { return 2 },
				func(x int) int { return x * 2 },
			},
			returnValue: []Argument{4},
			err:         nil,
		},
		{
			desc: "No return value",
			from: []Argument{4},
			funcs: []Function{
				func(x int) int { return x + 2 },
				func(x int) int { return x * 2 },
				func(int) {},
			},
			returnValue: []Argument{},
			err:         nil,
		},
		{
			desc: "With error",
			from: []Argument{},
			funcs: []Function{
				func() int { return 0 },
				func(x int) error { return errors.New("a generic error") },
			},
			returnValue: []Argument{},
			err:         errors.New("a generic error"),
		},
		{
			desc: "Without argument feedback",
			from: []Argument{},
			funcs: []Function{
				func() {},
				func() {},
			},
			returnValue: []Argument{},
			err:         nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rv, err := New().From(tc.from...).Chain(tc.funcs...).Unwrap()
			assert.Equal(t, tc.returnValue, rv)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestChainDummy(t *testing.T) {
	testCases := []struct {
		desc        string
		from        []Argument
		funcs       []Function
		returnValue []Argument
		err         error
	}{
		{
			desc: "Regular chain",
			from: []Argument{4},
			funcs: []Function{
				func(x int) int { return x + 2 },
				func(x int) int { return x * 2 },
			},
			returnValue: []Argument{12},
			err:         nil,
		},
		{
			desc: "With error in the end of the chain",
			from: []Argument{},
			funcs: []Function{
				func() int { return 2 },
				func(x int) (int, error) { return x, errors.New("a generic error") },
			},
			returnValue: []Argument{2, errors.New("a generic error")},
			err:         nil,
		},
		{
			desc: "With cascading error",
			from: []Argument{},
			funcs: []Function{
				func() error { return errors.New("a generic error") },
				func(e error) error { return e },
				func(e error) error { return e },
			},
			returnValue: []Argument{errors.New("a generic error")},
			err:         nil,
		},
		{
			desc: "With cascading error",
			from: []Argument{errors.New("a generic error")},
			funcs: []Function{
				func(e error) error { return e },
				func(e error) error { return e },
			},
			returnValue: []Argument{errors.New("a generic error")},
			err:         nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			rv, err := New().From(tc.from...).ChainDummy(tc.funcs...).Unwrap()
			assert.Equal(t, tc.returnValue, rv)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestNotAFunction(t *testing.T) {
	var err error

	errMsg := "(Error on Link: 0) Element isn't a function"

	// string
	_, err = New().Chain("not a function").Unwrap()
	assert.EqualError(t, err, errMsg)

	// number
	_, err = New().Chain(123).Unwrap()
	assert.EqualError(t, err, errMsg)

	// nil
	_, err = New().Chain(nil).Unwrap()
	assert.EqualError(t, err, errMsg)
}

func TestArgumentMismatch(t *testing.T) {
	f1 := func() (int, error) { return 0, nil }
	f2 := func(a, b int) (int, error) { return 0, nil }

	_, err := New().Chain(f1, f2).Unwrap()
	assert.EqualError(t, err, "(Error on Link: 1) 1 arg(s) provided, but function arity is 2")
}

func TestReset(t *testing.T) {
	returnOne := func() int { return 1 }

	res, err := New().Chain(returnOne).Reset().Unwrap()

	assert.Equal(t, []Argument{}, res)
	assert.NoError(t, err)
}
