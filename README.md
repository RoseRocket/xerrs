# xerrs [![Build Status](https://travis-ci.org/RoseRocket/xerrs.svg?branch=master)](https://travis-ci.org/RoseRocket/xerrs)

Extended Errors for Go

## What? Why?

It is extremely important to be able to effectively debug problems in production. However, depending
on the simplicty of your language's error handling, managing errors and debugging can become tedious
and difficult. For instance, GoLang provides a basic error which contains a simple value - but it
lacks verbose information about the cause of the error or the state of the stack at error-time.

xerrs uses original GoLang built-in error types but stores extra data.

These are the extra extended features which could be useful for debugging a big project:

1. **Cause** - the original error which happened in the system. This is the original error which you
   want to preserve
2. **Mask** - a mask error. It might be used whenever you want to mask some errors happening on
   production with a more generic error. Mask is very useful if you want to prevent original error
   return back to the client. Calling Error() will always return Mask if set.
3. **Stack** - a detailed snapshot of the execution stack at error-time. Helpful for debugging.
4. **Data** - a map which could be used to store custom data associated with an error.
5. **Wrap** and **Wrapf** can be used to create an annotated error chain, but
   take a lower precedence than **Mask**.

## Quick Usage

### Basics

```go
import "github.com/roserocket/xerrs"

//....

if data, err := MyFunc(); err != nil {
    err := xerrs.Extend(err) // extend error

    //....

    // Details returns cause error, mask if specified, and the stack (accepting the maximum stack height as parameter)
    // In this example only 5 last stack function calls will be printed out
    fmt.Println(xerrs.Details(err, 5))

    //....
}
```

### Deferred logging + masking example

```go
func DoSomething(w http.ResponseWriter, r *http.Request) {
    var err error

    const maxCallstack = 5 // only 5 last stack function calls will be printed out

    defer func() {
        // Details returns cause error, mask if specified, and the stack (accepting the maximum stack height as parameter)
        // In this example only 5 calls in the stack will be printed out
        fmt.Println(xerrs.Details(err, maxCallstack))
    }()

    var someModel Model
    if err = ReadJSONFromReader(r, &someModel); err != nil {
        err = xerrs.Extend(err)
        DoSomethingWithError(w, err.Error()) // Calling Error() without setting a mask will return the original error.
        return
    }

    if _, err = DBCreateMyModel(&someModel); err != nil {
        err = xerrs.Mask(err, errors.New("We are experiencing technical difficulties"))
        DoSomethingWithError(w, err.Error()) // Error() will return the masked error in this case.
        return
    }

    OutputDataToClient(w, &someModel)
}
```

### Custom data in error

```go
func DoSomething(w http.ResponseWriter, r *http.Request) {
    var err error

    //......

    if err = ReadJSONFromReader(r, someModel); err != nil {
        err = xerrs.Extend(err)

        xerrs.SetData(err, "some_key", "VALUE") // set custom error value
    }

    //......

    fmt.Println(xerrs.Details(err))
    fmt.Println(xerrs.GetData(err, "some_key")) // print custom error value

    //......
}
```

### Compare errors

```go
func VeryComplexLongFunction(arg1, arg2) error {
    var err error
    badErr := errors.New("EPIC FAIL")

    //......

    // convert error to an extended one and use it to for debugging purposes
    err = xerrs.Extend(err)

    //......

    if xerrs.IsEqual(err, badErr) {
        // errors are equal. We need to do something here
    }

    //......
}
```

## Docs

#### func New

```go
func New(string) error
```

New creates a new xerr with a supplied message

Note it will also set the stack

#### func Errorf

```go
func Errorf(string, ...interface{}) error
```

Errorf creates a new xerr based on a formatted message

Note it will also set the stack

#### func Extend

```go
func Extend(error) error
```

Extend creates a new xerr based on a supplied error

Note if error is nil then nil is returned

Note it will also set the stack

#### func Mask

```go
func Mask(error, error) error
```

Mask creates a new xerr based on a supplied error but also sets the mask error as well. Only mask
error value is returned back when Error() is called

Note if error is nil then nil is returned

Note if error is xerr then its mask value is updated

Note it will also set the stack

#### func IsEqual

```go
func IsEqual(error, error) bool
```

IsEqual is a helper function to compare if two errors are equal

Note if one of those errors is xerr then it's Cause is used for comparison

#### func Cause

```go
func Cause(error) error
```

Cause returns xerr's cause error

Note if error is not xerr then argument error is returned back

#### func SetData

```go
func SetData(error, string, interface{})
```

SetData sets custom data stored in xerr

Note if error is not xerr then function does not do anything

#### func GetData

```go
func GetData(error, string) (interface{}, bool)
```

GetData returns custom data stored in xerr

Note if error is not xerr then (nil, false) is returned

#### func Stack

```go
func Stack(error) []StackLocation
```

Stack returns stack location array

Note if error is not xerr then nil is returned

#### func Details

```go
func Details(error, int) string
```

Details returns a printable string which contains error, mask and stack

Note maxStack can be supplied to change number of printer stack rows

Note if error is not xerr then Error() is returned

## What are the alternatives?

xerrs library was partially inspired by [juju/errors](https://github.com/juju/errors)

[pkg/errors](https://github.com/pkg/errors)

Also there are
[new ideas and drafts for Go error handling](https://go.googlesource.com/proposal/+/master/design/go2draft.md)
which might change the way error is being handled in the future.

## LICENSE

see [LICENSE](./LICENSE)
