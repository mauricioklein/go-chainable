package chainable

import (
	"errors"
	"reflect"
)

var (
	errNotAFunction     = errors.New("fn is not a function")
	errArgumentMismatch = errors.New("fn argument mismatch")
)

// Chainable defines a new chainable structure
type Chainable struct {
	from  []interface{}
	links []Link
}

// Link represents a function in a chain
type Link struct {
	fn          interface{}
	handleError bool
}

// New instantiate a new chainable instance
func New() *Chainable {
	return &Chainable{}
}

// From defines the arguments provided to the first
// link in a chain
func (c *Chainable) From(args ...interface{}) *Chainable {
	c.from = args
	return c
}

// Chain appends a new function to the chain, with error
// checking enabled
func (c *Chainable) Chain(funcs ...interface{}) *Chainable {
	for _, fn := range funcs {
		c.addFunc(fn, true)
	}
	return c
}

// ChainDummy appends a new function to the chain, with error
// checking disabled
func (c *Chainable) ChainDummy(funcs ...interface{}) *Chainable {
	for _, fn := range funcs {
		c.addFunc(fn, false)
	}
	return c
}

// Unwrap process a chain, returning the result and error
// of the last execution (success or error)
func (c *Chainable) Unwrap() ([]interface{}, error) {
	v := c.from
	var err error

	for _, link := range c.links {
		if v, err = processLink(link, v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

// addFunc add a new function to the chain, creating the underlying link
func (c *Chainable) addFunc(fn interface{}, handleError bool) {
	c.links = append(c.links, Link{
		fn:          fn,
		handleError: handleError,
	})
}

// processLink calls the function fn associated to the link, transforming args using reflection
func processLink(link Link, args []interface{}) ([]interface{}, error) {
	vfn := reflect.ValueOf(link.fn)
	vfnType := vfn.Type()

	// check if it's a function
	if vfnType.Kind() != reflect.Func {
		return nil, errNotAFunction
	}

	// check if args matches the function arity
	if !vfnType.IsVariadic() && (vfnType.NumIn() != len(args)) {
		return nil, errArgumentMismatch
	}

	// call the function
	out := []interface{}{}
	for _, o := range vfn.Call(reflectArgs(args)) {
		out = append(out, o.Interface())
	}

	// if the last returned value for the function
	// is an error, cast the error and return it
	// along with the rest of values
	if link.handleError && doesReturnError(vfnType) {
		err, _ := out[len(out)-1].(error)
		return out[:len(out)-1], err
	}

	return out, nil
}

// reflectArgs transforms the args list in a list of
// reflect.Value, used to call a function using reflection
func reflectArgs(args []interface{}) []reflect.Value {
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
