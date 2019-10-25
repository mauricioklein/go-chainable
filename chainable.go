package chainable

import (
	"reflect"
)

// Function defines a function added to the chain
type Function interface{}

// Argument defines an argument provided to the functions in a chain
type Argument interface{}

// Chainable defines the chain logic.
type Chainable struct {
	from  []Argument
	links []link
}

// link represents a function in a chain
type link struct {
	fn          Function
	handleError bool
}

// New instantiate a new Chainable instance
func New() *Chainable {
	return &Chainable{}
}

// From defines the arguments provided to the first
// function in a chain
func (c *Chainable) From(args ...Argument) *Chainable {
	c.from = args
	return c
}

// Chain appends a new function to the chain, with error
// handling enabled
func (c *Chainable) Chain(funcs ...Function) *Chainable {
	c.chainFuncs(funcs, true)
	return c
}

// ChainDummy appends a new function to the chain, with error
// handling disabled
func (c *Chainable) ChainDummy(funcs ...Function) *Chainable {
	c.chainFuncs(funcs, false)
	return c
}

// Unwrap process the chain.
// It returns a []Argument with the arguments returned by the last function in the chain (or nil, if the chain is broken)
// and an error, returned by any of the functions in the chain with error-handling enabled (or nil, if success)
func (c *Chainable) Unwrap() ([]Argument, error) {
	v := c.from
	var err error

	for linkIndex, link := range c.links {
		if v, err = link.process(linkIndex, v); err != nil {
			return v, err
		}
	}

	return v, nil
}

// Reset cleanups a chain, removing all
// the links and initial values
func (c *Chainable) Reset() *Chainable {
	c.from = []Argument{}
	c.links = []link{}

	return c
}

// chainFuncs add new functions to the chain, creating the underlying link
func (c *Chainable) chainFuncs(funcs []Function, handleError bool) {
	for _, fn := range funcs {
		c.links = append(c.links, link{
			fn:          fn,
			handleError: handleError,
		})
	}
}

// process calls the function fn associated to the link, transforming args using reflection
func (lk *link) process(linkIndex int, args []Argument) ([]Argument, error) {
	vfn := reflect.ValueOf(lk.fn)
	if err := validateFunc(linkIndex, vfn); err != nil {
		return nil, err
	}

	vfnType := vfn.Type()
	if err := validateArgs(linkIndex, vfnType, args); err != nil {
		return nil, err
	}

	// call the function
	out := []Argument{}
	for _, o := range vfn.Call(reflectArgs(vfnType, args)) {
		out = append(out, o.Interface())
	}

	// if the last returned value for the function
	// is an error, cast the error and return it
	// along with the rest of values
	if lk.handleError && doesReturnError(vfnType) {
		err, _ := out[len(out)-1].(error)
		return out[:len(out)-1], err
	}

	return out, nil
}

// validateFunc validates if obj is a function
func validateFunc(linkIndex int, obj reflect.Value) error {
	// Zero value reflected: not a valid function
	if !isFunc(obj) {
		return notAFunctionError(linkIndex)
	}

	return nil
}

// isFunc returns a boolean indicating if obj is a function object
func isFunc(obj reflect.Value) bool {
	// Zero value reflected: not a valid function
	if obj == (reflect.Value{}) {
		return false
	}

	if obj.Type().Kind() != reflect.Func {
		return false
	}

	return true
}

// validateFuncArgs validates if len(args) matches the arity of fn
func validateArgs(linkIndex int, fn reflect.Type, args []Argument) error {
	if !fn.IsVariadic() && (fn.NumIn() != len(args)) {
		return argumentMismatchError(linkIndex, len(args), fn.NumIn())
	}

	return nil
}

// reflectArgs transforms the args list in a list of
// reflect.Value, used to call a function using reflection
func reflectArgs(fnType reflect.Type, args []Argument) []reflect.Value {
	in := make([]reflect.Value, len(args))

	for k, arg := range args {
		if arg == nil {
			// Use the zero value of the function parameter type,
			// since "reflect.Call" doesn't accept "nil" parameters
			in[k] = reflect.New(fnType.In(k)).Elem()
		} else {
			in[k] = reflect.ValueOf(arg)
		}
	}

	return in
}

// doesReturnError returns true if the last return
// value for vfn is a type error
func doesReturnError(vfnType reflect.Type) bool {
	if vfnType.NumOut() == 0 {
		return false
	}

	return vfnType.Out(vfnType.NumOut()-1) == reflect.TypeOf((*error)(nil)).Elem()
}
