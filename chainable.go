package main

import (
	"fmt"
)

// Args defines a list of arguments accepted
// and returned by chainable functions
type Args []interface{}

// NoArgs defines an empty list of arguments
var NoArgs = Args{}

// ChainableFunc defines the function signature accepted
// by Chainable
type ChainableFunc func(Args) (Args, error)

// Chainable defines the chainable structure
type Chainable struct {
	initialValue Args
	calls        []ChainableFunc
}

// New initalizes a chainable structure
func New() *Chainable {
	return &Chainable{}
}

// InitialValue defines the value used as argument for the
// first function in the chain
func (c *Chainable) InitialValue(al Args) *Chainable {
	c.initialValue = al
	return c
}

// Chain adds a new function to the chain
func (c *Chainable) Chain(f ChainableFunc) *Chainable {
	c.calls = append(c.calls, f)
	return c
}

// Wrap process a chain and returns the last return
// value and error
func (c *Chainable) Wrap() (Args, error) {
	lastRet := c.initialValue
	var err error

	for _, c := range c.calls {
		if lastRet, err = c(lastRet); err != nil {
			return nil, err
		}
	}

	return lastRet, nil
}

func main() {
	chainable := New()

	f1 := func(Args) (Args, error) { return NoArgs, nil }
	f2 := func(Args) (Args, error) { return NoArgs, nil }
	f3 := func(Args) (Args, error) { return NoArgs, nil }

	ret, err := chainable.
		Chain(f1).
		Chain(f2).
		Chain(f3).
		Wrap()

	fmt.Printf("Ret, Err: %+v, %+v\n", ret, err)

}
