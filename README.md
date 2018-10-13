# xerrs

Extended Errors for Go

## What? Why?

It is extremely important to be able to effectively debug problems in production. However, depending
on the simplicty of your language's error handling, managing errors and debugging can become tedious
and difficult. For instance, GoLang provides a basic error which contains a simple value - but it
lacks verbose information about the cause of the error or the state of the stack at error-time.

One solution might be to create a new struct which mimics errors but add extra data. However, this
solution presents problems for refactoring and future error management, especially if your code base
is large. Moreover, a custom error object might prevent the easy use of third-party Go libraries,
most of which use native Go errors.

XErrs uses original GoLang built-in error types and stores all extra data as a string within the
error itself.

These are the extra extended features which could be extremly useful for any debugging or big
project:

1. **CauseError** - the original error which happened in the system. This is the original built-in
   GoLang error which you want to preserve
2. **MaskError** - a client-facing error. It might be used whenever you want to mask some errors
   happening on production with a more generic or semantic error
3. **Stack** - a detailed snapshot of the execution stack at error-time. It might be used to help
   you to debug errors
4. **Data** - a mao which contains extra custom data. It could be used for any other data you might
   want to associate with an error.

## How it works

To provide added added context to an error, we store an original error object with extra data as a
serialized object inside a new GoLang error. Using xerrs functions we can extract any relevant data
from this new error or convert it back to the original error schema.

## Bonus Points!

Since the error message is a JSON string it makes very easy to share error objects between systems
written in different programming languages. For instance, extended error values could be shared by a
system implemented in Python. Alternatively, your Go code could create an error and then pass it to
any other systems, regardless of their implementation.

## Quick Usage

### Basics

```go
import "github.com/roserocket/xerrs"

//....

if data, err := MyFunc(); err != nil {
    err = xerrs.Extend(err) // extend error

    //....

    xerrs.Stack(err) // print detailed error information including its stack
}
```

### Deferred logging + masking example

```go
func DoSomething(w http.ResponseWriter, r *http.Request) {
    var err error

    defer func() {
        xerrs.Stack(err)
    }()

    someModel := &Model{}
    err = ReadJSONFromReader(r, someModel)
    if err != nil {
        err = xerrs.Extend(err)
        DoSomethingWithError(w, xerrs.Error(err))
        return
    }

    _, err = DBCreateMyModel(someModel)
    if err != nil {
        err = xerrs.MaskError(err, errors.New("We are experiencing technical difficulties"))
        DoSomethingWithError(w, xerrs.Error(err))
        return
    }

    OutputDataToClient(w, &someModel)
}
```

### Backward compatibility with existing code

```go
func VeryComplexLongFunction(arg1, arg2) error {
    var err error
    // convert error to an extended one and use it to for debugging purposes

    if err {
        err = xerrs.Extend(err)
    }

    //......

    if err {
        err = xerrs.Extend(err)
    }

    //......
    xerrs.Stack(err)

    //......

    // convert back to the regular GoLang error type so that other function
    // will work without a single change
    // If err is nil then nil is returned back
    return xerrs.ToError(err)
}
```

### Preserving error stack

```go
func innerFunc() error {
    //......

    if err = something; err != nil {
        return err
    }

    //......

    if err = somethingElse; err != nil {
        return xerrs.Extend(err)
    }

    //......
}

func outterFunc() {
    //......

    if err = innerFunc(); err != nil {
        // At this point if returned err from innerFunc() is XErr then it stays
        // unchanged with a preserved error stack starting at innerFunc
        // If err is just a regular error then it will be converted into XErr
        // where stack will be starting from the outterFunc()
        err = xerrs.Extend(err)

        // do something here with err
    }

    //......
}
```

## Docs

#### func Extend

```go
func Extend(error) error
```

Extend takes the original error and returns an extended one. Its value is equal to the stringified
JSON of the Extended Error object containing the original error. Exerr's functions work primarily on
these extended errors.

Note that the original argument is returned if error is already an extended error.

#### func MaskError

```go
func MaskError(error, error) error
```

MaskError works identically as `Extend()` however the second argument is mask error. Mask is used to
conceal the real error which happened in the system. This could be useful if you need to preserve
the original error without exposing it to the client.

Note that mask will be changed if the first argument is an extended error.

#### func Cause

```go
func Cause(error) error
```

Cause will return an original error which was extended.

Note that the first argument is returned back if it is not an extended error.

#### func SetData

```go
func SetData(error, string, interface{}) error
```

SetData sets custom Data Error property of an extended error. Could be handy if you need to pass any
extra values with an error.

Note that the first argument is returned back if it is not an extended error.

#### func Data

```go
func GetData(error, string) (string, bool)
```

GetData will return custom Data Error property by name of an extended error.

Note that the second return vaue is false if error is not an extended one.

#### func Error

```go
func Error(error) string
```

Error is the same as `err.Error()` call. However if the argument is an extended error then `Error()`
will be called on the mask error. This function hides access to the original Cause error and could
be used for sending an error message to the client.

#### func Trace

```go
func Stack(error) string
```

Stack returns a detailed string for logging and debugging the error. This detailed string will
consist of the original error, its mask, and log-ready stack lines.

Note that the original argument is returned if error is already an extended error.

#### func IsEqual

```go
func IsEqual(error, error) bool
```

Returns true if two supplied errors have the same value.

Note If one of those errors is XErr then its Cause value will be used for comparing.

#### func ToXErr

```go
func ToXErr(error) (*XErr, bool)
```

ToXErr takes an extended error and converts it to the `XErr` object.

Note that the second return value is false if error is not an extended one.

## Possible problems to be aware of

Developers should be aware that extended error is just serialized json string and you might want to
convert it back to the original error at some point in the lifespan of your app. If you start
comparing extended error to other errors this comparison will most likely fail. Returning this
serialized string is not really a client friendly error, plus you would expose your code stack to
the client... which is not good.

Also you cannot use basic comparison == between the original error and xerrs.Cause(err) due to
marshalling and unmarshalling logic under the hood. Instead xerrs.IsEqual() should be used.

Example:

```go
func ExampleFunc() {
    var err error
    var fault = errors.New("Something very bad")

    err = xerrs.Extend(fault)

    //......

    if xerrs.Cause(err) == fault {
        // this would never work. Even though fault is used as a cause for xerrs
        // it will not be the same after it is unmarshalled
    }

    if xerrs.IsEqual(err, fault) {
        // this would work just fine
    }
```

## What are the alternatives?

xerrs library was partially inspired by [juju/errors](https://github.com/juju/errors)

[pkg/errors](https://github.com/pkg/errors)

Also there are
[new ideas and drafts for Go error handling](https://go.googlesource.com/proposal/+/master/design/go2draft.md)
which might change the way error is being handled in the future.

## LICENSE

see [LICENSE](./LICENSE)
