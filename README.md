# xerrs

Extended Errors for Go

## What? Why?

It is extremely important being able effectively debug problems on production. However it could be
very tedious and hard to manage especially if your language provides very simple tools to work with.
Just like GoLang provides a basic error which contains a simple value.

One of the solution would be creating new struct which mimics error but add extra data. However it
could be hard to manage or even refactor if your existing code base consist many lines. The fact
that most of the third-party libraries use built-in go error types would not make things simpler by
forcing constant types conversion or writing special wrappers.

XErrs uses original GoLang built-in error types but would store all extra data as a string within
the error itself.

These are the extra extended features which could be extremly useful for any debugging or big
project:

1. **CauseError** - an original error which happened in the system. This is the origin basic GoLang
   error which you want to preserve
2. **MaskError** - a client-facing error. It might be used when you want to mask some errors
   happening on production with a more generic error returned back to whatever client consuming your
   code.
3. **Stack** - a detailed execution stack at the time when error was generated. It might be used to
   help you to debug your errors
4. **Data** - a string which contains any extra custom data stored in the error. It might be used if
   you want to have custom data, error codes or whatever you want serialized as json string to be
   passed along with the initial error

## How it works

We serialize an object as a string which contains an original error along with lots of other extra
data and store it inside the GoLang error. Using xerrs functions we can extract any relevant data or
even convert it back to the original error.

## Bonus Points!

Since the error message is a JSON string it makes way easier to share error object between systems
written in different programming languages. Extended error value could be created by (let's say
system written in Python) and passed into your Go one. Or vice versa your Go code could create an
error and then pass it to whatever other system writte on any other language to process.

## Quick Usage

### Basics

```go
import "github.com/roserocket/xerrs"

//....

if data, err := MyFunc(); err != nil {
    err = xerrs.ExtendError(err) // extend error

    //....

    xerrs.Stack(err) // print detailed error information including its stack
}
```

### Deferred logging + masking

```go
func CreateDock(w http.ResponseWriter, r *http.Request) {
    var err error

    defer func() {
        xerrs.Stack(err)
    }()

    dock := &models.Dock{}
    err = utils.ReadJSON(r, dock)
    if err != nil {
        err = xerrs.ExtendError(err)
        WriteJSONError(w, xerrs.Error(err), http.StatusBadRequest)
        return
    }

    dock, err = db.DocksCreate(dock)
    if err != nil {
        err = xerrs.MaskError(err, errors.New("We are experiencing technical difficulties"))

        WriteJSONError(w, xerrs.Error(err), http.StatusInternalServerError)
        return
    }

    WriteJSON(w, &utils.JSONResponse{Data: dock})
}
```

### Backward compatibility with existing code

```go
func VeryComplexLongFunction(arg1, arg2) error {
    var err error
    // convert error to an extended one and use it to for debugging purposes

    if err {
        err = xerrs.ExtendError(err)
    }

    //......

    if err {
        err = xerrs.ExtendError(err)
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

## Docs

#### func ExtendErrory

```go
func ExtendError(error) error
```

ExtendError takes original error and returns an extended one. Its value is equal to stringified JSON
of the Extended Error object containing an original error. Most of this library functions would work
only on such extended errors.

Note that the original argument is returned if error is already an extended error.

#### func MaskError

```go
func MaskError(error, error) error
```

MaskError works identically as `ExtendError()` however the second argument is mask error. Mask is
used to conceal the real error which happened in the system. Could be handy if you do not want to
expose to the client a real error which has happend but at the same time need to preserve the
original error for logging purposes.

Note that mask will be changed if the first argument is an extended error.

#### func Cause

```go
func Cause(error) error
```

Cause will return an original error which was extended.

Note that the first argument is returned back if it is not an extended error.

#### func SetData

```go
func SetData(error, string) error
```

SetData sets Data property of an extended error. Could be handy if you need to pass along any extra
data or a serialized as a string object.

Note that the first argument is returned back if it is not an extended error.

#### func Data

```go
func Data(error) string, bool
```

Data will return data string property extended error.

Note that the second return vaue is false if error is not an extended one.

#### func Error

```go
func Error(error) string
```

Error is the same as `err.Error()` call. However if the argument is an extended error then `Error()`
will be called on the mask error. This function hides access to the original Cause error and could
be used for sending to the client side.

#### func Trace

```go
func Stack(error) string
```

Stack returns a full detailed string ready for logging and helping to debug. This detailed string
will consist o the original error, its mask as well as log-ready stack lines.

Note that the original argument is returned if error is already an extended error.

#### func ToXErr

```go
func ToXErr(error) (*XErr, bool)
```

ToXErr takes an extended error and converts it to the `XErr` object.

Note that the second return value is false if error is not an extended one.

## Possible problems to be aware of

Developers should be aware that extended error is just serialized json string and you might want to
convert it back to the original error at some point in the lifespan of your app. If you start
comparing extended error to other errors this comparison most likely would fail. Returning this
serialized string is not really a client friendly error, plus you would expose your code stack to
the client... which is not good.

## What are the alternatives?

xerrs library was partially inspired by [juju/errors](https://github.com/juju/errors)

[juju/errgo](https://github.com/juju/errgo)

Also there are
[new ideas and drafts for Go error handling](https://go.googlesource.com/proposal/+/master/design/go2draft.md)
which might change the way error is being handled in the future.

## LICENSE

see [LICENSE](./LICENSE)
