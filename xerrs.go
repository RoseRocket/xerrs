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
		data:  make(map[string]interface{}),
		cause: errors.New(message),
		mask:  nil,
		stack: getStack(stackFunctionOffset),
	}
}

// Errorf - creates a new xerr based on a formatted message.
// It will also set the stack.
func Errorf(format string, args ...interface{}) error {
	return &xerr{
		data:  make(map[string]interface{}),
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
		data:  make(map[string]interface{}),
		cause: err,
		mask:  nil,
		stack: getStack(stackFunctionOffset),
	}
}

// Mask - creates a new xerr based on a supplied error but also sets the mask error as well
// When Error() is called on the error only mask error value is returned back
// If err is nil then nil is returned
// It will also set the stack.
func Mask(err, mask error) error {
	if err == nil {
		return nil
	}

	return &xerr{
		data:  make(map[string]interface{}),
		cause: err,
		mask:  mask,
		stack: getStack(stackFunctionOffset),
	}
}

// IsEqual - helper function to compare if two erros are equal
// If one of those errors are xerr then its Cause is used for comparison
func IsEqual(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	if err1 == nil || err2 == nil {
		return false
	}

	x1, ok1 := err1.(*xerr)
	x2, ok2 := err2.(*xerr)

	if ok1 && ok2 {
		return x1.Cause().Error() == x2.Cause().Error()
	} else if !ok1 && ok2 {
		return err1.Error() == x2.Cause().Error()
	} else if ok1 && !ok2 {
		return x1.Cause().Error() == err2.Error()
	} else {
		return err1.Error() == err2.Error()
	}
}

// Cause - returns xerr cause error
func (x *xerr) Cause() error {
	return x.cause
}

func (x *xerr) Error() string {
	if x.mask == nil {
		return x.cause.Error()
	}

	return x.mask.Error()
}

// Mask - sets the mask in xerr
func (x *xerr) Mask(err error) {
	x.mask = err
}

// GetData - returns custom data stored in xerr
func (x *xerr) GetData(name string) (value interface{}, ok bool) {
	value, ok = x.data[name]
	if !ok {
		return
	}

	return
}

// SetData - sets custom data stored in xerr
func (x *xerr) SetData(name string, value interface{}) {
	x.data[name] = value
}

// Stack - returns stack location array
func (x *xerr) Stack() []StackLocation {
	return x.stack
}

// Details - returns a printable string which contains error, mask and stack
// maxStack can be supplied to change number of printer stack rows
func (x *xerr) Details(maxStack int) string {
	const newLine = "\n"

	result := []string{""}

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
