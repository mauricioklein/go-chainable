package chainable

import (
	"errors"
	"reflect"
)

var (
	errNotAFunction             = errors.New("fn is not a function")
	errArgumentMismatch         = errors.New("fn argument mismatch")
	errRightmostValueNotAnError = errors.New("rightmost return value should be an error")
)

// Chainable defines a new chain structure
type Chainable struct {
	initialValue []interface{}
	functions    []interface{}
}

// New instantiate a new chain instance
func New() *Chainable {
	return &Chainable{}
}

// From defines the arguments provided to the first
// function in a chain
func (c *Chainable) From(args []interface{}) *Chainable {
	c.initialValue = args
	return c
}

// Chain adds a new function to the end of the chain
func (c *Chainable) Chain(fn interface{}) *Chainable {
	c.functions = append(c.functions, fn)
	return c
}

// Unwrap process a chain, returning the result and error
// of the last execution (success or error)
func (c *Chainable) Unwrap() ([]interface{}, error) {
	v := c.initialValue
	var err error

	for _, fn := range c.functions {
		v, err = callFn(fn, v)
		if err != nil {
			return nil, err
		}

	}

	return v, nil
}

// callFn calls the function fn, transforming args
// using reflection
func callFn(fn interface{}, args []interface{}) ([]interface{}, error) {
	vfn := reflect.ValueOf(fn)

	// check if it's a function
	if vfn.Type().Kind() != reflect.Func {
		return nil, errNotAFunction
	}

	// build the reflected args
	reflectedArgs, err := buildReflectedArgs(vfn, args)
	if err != nil {
		return nil, err
	}

	// call the function
	out := []interface{}{}
	for _, o := range vfn.Call(reflectedArgs) {
		out = append(out, o.Interface())
	}

	err, ok := out[len(out)-1].(error)
	if !ok && out[len(out)-1] != nil {
		return nil, errRightmostValueNotAnError
	}

	return out[:len(out)-1], err
}

// buildReflectedArgs transforms the args list in a list of
// reflect.Value, used to call a function using reflection
func buildReflectedArgs(vfn reflect.Value, args []interface{}) ([]reflect.Value, error) {
	if !vfn.Type().IsVariadic() && (vfn.Type().NumIn() != len(args)) {
		return nil, errArgumentMismatch
	}

	in := make([]reflect.Value, len(args))
	for k, arg := range args {
		in[k] = reflect.ValueOf(arg)
	}

	return in, nil
}
