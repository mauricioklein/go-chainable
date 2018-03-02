package chainable

import (
	"errors"
	"reflect"
)

var (
	errNotAFunction     = errors.New("fn is not a function")
	errArgumentMismatch = errors.New("fn argument mismatch")
)

// Chain defines a new chain structure
type Chain struct {
	from  []interface{}
	funcs []interface{}
}

// New instantiate a new chain instance
func New() *Chain {
	return &Chain{}
}

// From defines the arguments provided to the first
// link in a chain
func (c *Chain) From(args ...interface{}) *Chain {
	c.from = args
	return c
}

// Link appends a new function to the chain
func (c *Chain) Link(funcs ...interface{}) *Chain {
	for _, fn := range funcs {
		c.funcs = append(c.funcs, fn)
	}

	return c
}

// Unwrap process a chain, returning the result and error
// of the last execution (success or error)
func (c *Chain) Unwrap() ([]interface{}, error) {
	v := c.from
	var err error

	for _, fn := range c.funcs {
		if v, err = callFn(fn, v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

// callFn calls the function fn, transforming args
// using reflection
func callFn(fn interface{}, args []interface{}) ([]interface{}, error) {
	vfn := reflect.ValueOf(fn)
	vfnType := vfn.Type()

	// check if it's a function
	if vfn.Type().Kind() != reflect.Func {
		return nil, errNotAFunction
	}

	// check if args matches the function arity
	if !vfnType.IsVariadic() && (vfnType.NumIn() != len(args)) {
		return nil, errArgumentMismatch
	}

	// build the reflected args
	reflectedArgs := buildReflectedArgs(args)

	// call the function
	out := []interface{}{}
	for _, o := range vfn.Call(reflectedArgs) {
		out = append(out, o.Interface())
	}

	// if the last returned value for the function
	// is an error, cast the error and return it
	// along with the rest
	if doesReturnError(vfnType) {
		err, _ := out[len(out)-1].(error)
		return out[:len(out)-1], err
	}

	return out, nil
}

// buildReflectedArgs transforms the args list in a list of
// reflect.Value, used to call a function using reflection
func buildReflectedArgs(args []interface{}) []reflect.Value {
	in := make([]reflect.Value, len(args))

	for k, arg := range args {
		in[k] = reflect.ValueOf(arg)
	}

	return in
}

// doesReturnError returns true if the last return
// value for vfn is a type error
func doesReturnError(vfnType reflect.Type) bool {
	return vfnType.Out(vfnType.NumOut()-1) == reflect.TypeOf((*error)(nil)).Elem()
}
