# Go-chainable

Go chain functions execution, with support to feedback and error-handling.

## Motivation

Found in Elixir, the pipe operator allows to chain 2 or more functions, feedbacking the return value of the previous
as the input of the next one.

This library provides a similar functionality, given Golang limitations.

## Installation

```go
go get github.com/mauricioklein/go-chainable
```

## Usage example

```go
import chainable "github.com/mauricioklein/go-chainable"

sum2 := func(x int) int { return x + 2; }
mul2 := func(x int) int { return x * 2; }

chainable.New().
    Link(
        sum2, 
        mul2,
    ).
    Unwrap()
```

Also, the library supports error handling. So, in case one of the links returns an error,
the chain is broken and the error code is returned. The error must be the last returned
value, otherwise, it's treated as a regular return value.

```go
import chainable "github.com/mauricioklein/go-chainable"

sum2 := func(x int) (int, error) { return 0, errors.New("a random error") }
mul2 := func(x int) int { return x * 2; }

chainable.New().
    Link(
        sum2, 
        mul2,
    ).
    Unwrap()
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
