package chainable

import (
	"fmt"
)

// notAFunctionError returns the formatted error
// when an element in the chain isn't a function
func notAFunctionError(linkIndex int) error {
	return fmt.Errorf("(Error on Link: %d) Element isn't a function", linkIndex)
}

// argumentMismatchError returns the formatted error
// when the number of arguments provided doesn't match
// the function's arity
func argumentMismatchError(linkIndex, nArgs, fnArity int) error {
	return fmt.Errorf("(Error on Link: %d) %d arg(s) provided, but function arity is %d", linkIndex, nArgs, fnArity)
}
