# Go-chainable

[![Build Status](https://travis-ci.org/mauricioklein/go-chainable.svg?branch=master)](https://travis-ci.org/mauricioklein/go-chainable)

Chain functions execution, with support to argument's feedback and error-handling.

## Motivation

Found in Elixir, the pipe operator allows us to chain 2 or more functions, feedbacking the return value of the previous
as the input of the next one.

This library provides a similar functionality to Golang, given language limitations.

## Installation

```go
go get github.com/mauricioklein/go-chainable
```

### Usage

Go-chainable chains the result of the previous function as the input of the next one.
The number of arguments returned by the previous function must match the arity
of the next one.

```go
// create a new chain
cn := chainable.New()

joinWithSpaces := func(strs []string) string { return strings.Join(strs, " ") }

cn.
    From([]string{"hello", "world"}). // "From" defines the arguments for the first method in the chain
    Chain(joinWithSpaces).            
    Chain(strings.Title)

// "Unwrap" processes the entire chain, returning the result of the last 
// function as an []interface{}, which can be casted to the real values
res, err := cn.Unwrap()

fmt.Printf("Result: %s\n", res[0].(string)) // "Hello World"
fmt.Printf("Error: %v\n", err) // nil

```

`Go-chainable` also supports error handling. In this case, if one of the functions in the chain
returns an error, the chain is broken and the error is returned.

To support error handling, the chained function must return an error as the last argument.

```go
raiseErr := func() error { return errors.New("a generic error") }
returnOne := func() int { return 1 }

result, err := cn.
    Chain(raiseErr).
    Chain(returnOne). // not called, since the last function broke the chain
    Unwrap()

fmt.Printf("Result: %d\n", result) // nil
fmt.Printf("Error: %+v\n", err) // "a generic error"
```

In case the error shouldn't be handled, but chained to the next function, you can use instead the method `ChainDummy`:

```go
raiseErr := func() error { return errors.New("a generic error") }
returnOneAndError := func(err error) (int, error) { return 1, err }

result, err := cn.
    ChainDummy(raiseErr). // doesn't handle the error, but chains as input for the next function
    Chain(returnOneAndError).
    Unwrap()

fmt.Printf("Result: %d\n", result) // []interface{1}
fmt.Printf("Error: %+v\n", err) // "a generic error"
```

Both `Chain` and `ChainDummy` are variadic functions, so you can build a chain:

```go
// separatelly...
cn.
    Chain(f1).
    Chain(f2).
    Chain(f3).
    Unwrap()

// ...or together
cn.
    Chain(f1, f2, f3).
    Unwrap()
```

Finally, to reset a chain and make it ready to be reused, just call the method `Reset`:

```go
cn.Reset()
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
