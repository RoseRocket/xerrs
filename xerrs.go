package xerrs

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// This value represents the offset in the stack array. We want to keep this
// number at 2 so that we do not see XErr functions in the stack
const stackFunctionOffset = 2

type xerr struct {
	data  map[string]interface{}
	cause error
	mask  error
	stack []StackLocation
}

func (x *xerr) Error() string {
	if x.mask == nil {
		return x.cause.Error()
	}

	return x.mask.Error()
}

// StackLocation - A helper struct function which represents one step in the execution stack
type StackLocation struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Returns a string which represents StackLocation. Used for logging.
func (location StackLocation) String() string {
	return fmt.Sprintf("%s [%s:%d]", location.Function, location.File, location.Line)
}

// New - creates a new xerr with a supplied message.
// It will also set the stack.
func New(message string) error {
	return &xerr{
		data:  nil,
		cause: errors.New(message),
		mask:  nil,
		stack: getStack(stackFunctionOffset),
	}
}

// Errorf - creates a new xerr based on a formatted message.
// It will also set the stack.
func Errorf(format string, args ...interface{}) error {
	return &xerr{
		data:  nil,
		cause: fmt.Errorf(format, args...),
		mask:  nil,
		stack: getStack(stackFunctionOffset),
	}
}

// Extend - creates a new xerr based on a supplied error.
// If err is nil then nil is returned
// It will also set the stack.
func Extend(err error) error {
	if err == nil {
		return nil
	}

	return &xerr{
		data:  nil,
		cause: err,
		mask:  nil,
		stack: getStack(stackFunctionOffset),
	}
}

// Mask - creates a new xerr based on a supplied error but also sets the mask error as well
// When Error() is called on the error only mask error value is returned back
// If err is nil then nil is returned
// If err is xerr then its mask value is updated
// It will also set the stack.
func Mask(err, mask error) error {
	if err == nil {
		return nil
	}

	if x, ok := err.(*xerr); ok {
		x.mask = mask
		return x
	}

	return &xerr{
		data:  nil,
		cause: err,
		mask:  mask,
		stack: getStack(stackFunctionOffset),
	}
}

// IsEqual - helper function to compare if two erros are equal
// If one of those errors are xerr then its Cause is used for comparison
func IsEqual(err1, err2 error) bool {
	cause1 := Cause(err1)
	cause2 := Cause(err2)

	if cause1 == nil && cause2 == nil {
		return true
	}

	if cause1 == nil || cause2 == nil {
		return false
	}

	return Cause(err1).Error() == Cause(err2).Error()
}

// Cause - returns xerr's cause error
// If err is not xerr then err is returned
func Cause(err error) error {
	if x, ok := err.(*xerr); ok {
		return x.cause
	}

	return err
}

// GetData - returns custom data stored in xerr
// If err is not xerr then (nil, false) is returned
func GetData(err error, name string) (value interface{}, ok bool) {
	var x *xerr

	x, ok = err.(*xerr)

	if !ok {
		return
	}

	if x.data == nil {
		ok = false
		return
	}

	value, ok = x.data[name]

	return
}

// SetData - sets custom data stored in xerr
// If err is not xerr then nothing happens
func SetData(err error, name string, value interface{}) {
	if x, ok := err.(*xerr); ok {
		if x.data == nil {
			x.data = make(map[string]interface{})
		}

		x.data[name] = value
	}
}

// Stack - returns stack location array
// If err is not xerr then nil is returned
func Stack(err error) []StackLocation {
	if x, ok := err.(*xerr); ok {
		return x.stack
	}

	return nil
}

// Details - returns a printable string which contains error, mask and stack
// maxStack can be supplied to change number of printer stack rows
// If err is not xerr then err.Error() is returned
func Details(err error, maxStack int) string {
	if err == nil {
		return ""
	}

	const newLine = "\n"

	result := []string{""}
	x, ok := err.(*xerr)

	if !ok {
		return err.Error()
	}

	result = append(result, fmt.Sprintf("[ERROR] %s", x.cause.Error()))
	if x.mask != nil && x.cause.Error() != x.mask.Error() {
		result = append(result, fmt.Sprintf("[MASK ERROR] %s", x.mask.Error()))
	}

	if len(x.stack) == 0 {
		return strings.Join(result, newLine)
	}

	result = append(result, "[STACK]:")

	top := maxStack
	if maxStack > len(x.stack) {
		top = len(x.stack)
	}

	for i := 0; i < top; i++ {
		result = append(result, x.stack[i].String())
	}

	return strings.Join(result, newLine)
}

// Returns execution Stack of the goroutine which called it in the form of StackLocation array
// skip - is a starting level on the execution stack where 0 = getStack() function itself, 1 = caller who called getStack(), and so forth
func getStack(skip int) []StackLocation {
	stack := []StackLocation{}

	i := 0
	for {
		pc, fn, line, ok := runtime.Caller(skip + i)
		if !ok {
			return stack
		}

		newStackLocation := StackLocation{
			Function: runtime.FuncForPC(pc).Name(),
			File:     fn,
			Line:     line,
		}

		stack = append(stack, newStackLocation)

		i++
	}
}
