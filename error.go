package chainable

import (
	"fmt"
)

// NotAFunctionError defines the error raised
// when a non-function object is chained
type NotAFunctionError struct {
	linkIndex int
}

// ArgumentMismatchError defines the error raised
// when the number of arguments provided to a function
// doesn't match the function arity
type ArgumentMismatchError struct {
	linkIndex int
	nArgs     int
	fnArity   int
}

// NewNotAFunctionError creates a new instance of NotAFunctionError
func NewNotAFunctionError(linkIndex int) *NotAFunctionError {
	return &NotAFunctionError{linkIndex}
}

// NewArgumentMismatchError creates a new instance of ArgumentMismatchError
func NewArgumentMismatchError(linkIndex, nArgs, fnArity int) *ArgumentMismatchError {
	return &ArgumentMismatchError{
		linkIndex: linkIndex,
		nArgs:     nArgs,
		fnArity:   fnArity,
	}
}

// Error provides the description of the error
func (err *NotAFunctionError) Error() string {
	return fmt.Sprintf("(Error on Link: %d) Element isn't a function", err.linkIndex)
}

// Error provides the description of the error
func (err *ArgumentMismatchError) Error() string {
	return fmt.Sprintf("(Error on Link: %d) %d arg(s) provided, but function arity is %d", err.linkIndex, err.nArgs, err.fnArity)
}
