# xerrs

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

## Quick Usage

### Basics

```go
import "github.com/roserocket/xerrs"

//....

if data, err := MyFunc(); err != nil {
    err := xerrs.Extend(err) // extend error

    //....

    if x, ok := err.(*xerr); ok {
        // In this example we only interested in the last 5 execution calls within the stack
        fmt.Println(x.Details(5)) // Details prints cause error, mask if specified, and stack (accepting the maximum stack height as parameter)
    } else {
        fmt.Println(err) // print basic error message if it is original error
    }

    //....
}
```

### Deferred logging + masking example

```go
func ErrPrintDetails(err error) {
    if x, ok := err.(*xerr); ok {
        fmt.Println(x.Details(5)) // Details prints cause error, mask if specified, and stack (accepting the maximum stack height as parameter)
    } else {
        fmt.Println(err) // print basic error message if it is original error
    }
}

func DoSomething(w http.ResponseWriter, r *http.Request) {
    var err error

    defer func() {
        ErrPrintDetails(err)
    }()

    someModel := &Model{}
    err = ReadJSONFromReader(r, someModel)
    if err != nil {
        err = xerrs.Extend(err)
        DoSomethingWithError(w, err.Error()) // Calling Error() without setting a mask will return the original error.
        return
    }

    _, err = DBCreateMyModel(someModel)
    if err != nil {
        err = xerrs.MaskError(err, errors.New("We are experiencing technical difficulties"))
        DoSomethingWithError(w, err.Error()) // Error() will return the masked error in this case.
        return
    }

    OutputDataToClient(w, &someModel)
}
```

### Custom data in error

```go
func ErrPrintDetails(err error) {
    if x, ok := err.(*xerr); ok {
        fmt.Println(x.Details(5))

        fmt.Println(x.GetData("VALUE")) // print custom error value
    } else {
        fmt.Println(err)
    }
}

func DoSomething(w http.ResponseWriter, r *http.Request) {

    //......

    var err error
    err = ReadJSONFromReader(r, someModel)
    if err != nil {
        err = xerrs.Extend(err)

        x, _ := err.(*xerr); ok {
            x.SetData("some_key", "VALUE")
        }
    }

    //......

    ErrPrintDetails(err)

    //......
}
```

### Compare errors

```go
func VeryComplexLongFunction(arg1, arg2) error {
    var err error
    badErr := errors.New("EPIC FAIL")
    // convert error to an extended one and use it to for debugging purposes

    if err {
        err = xerrs.Extend(err)
    }

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

Note if err is nil then nil is returned

Note it will also set the stack

#### func Mask

```go
func Mask(error, error) error
```

Mask creates a new xerr based on a supplied error but also sets the mask error as well When Error()
is called on the error only mask error value is returned back

Note if err is nil then nil is returned

Note it will also set the stack

#### func IsEqual

```go
func IsEqual(error, error) bool
```

IsEqual is a helper function compare if two erros are equal

Note if one of those errors are xerr then its Cause is used for comparison

#### xerr func Error

```go
func (x *xerr) Error() string
```

Error implements error interface Error() function

Note if xerr has a Mask error then Mask.Error() is returned back masking the original error

#### xerr func Cause

```go
func (x *xerr) Cause() error
```

Cause returns xerr cause error

#### xerr func Mask

```go
func (x *xerr) Mask(error)
```

Mask sets the mask in xerr

#### xerr func SetData

```go
func (x *xerr) SetData(string, interface{})
```

SetData sets custom data stored in xerr

#### xerr func GetData

```go
func (x *xerr) GetData(string) (interface{}, bool)
```

GetData returns custom data stored in xerr

#### xerr func Stack

```go
func (x *xerr) Stack() string
```

Stack returns stack location array

#### xerr func Details

```go
func (x *xerr) Details(int) string
```

Details returns a printable string which contains error, mask and stack

Note maxStack can be supplied to change number of printer stack rows

## What are the alternatives?

xerrs library was partially inspired by [juju/errors](https://github.com/juju/errors)

[pkg/errors](https://github.com/pkg/errors)

Also there are
[new ideas and drafts for Go error handling](https://go.googlesource.com/proposal/+/master/design/go2draft.md)
which might change the way error is being handled in the future.

## LICENSE

see [LICENSE](./LICENSE)
