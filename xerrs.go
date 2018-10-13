package xerrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Maximum number of function lines printed out for Stack() call
const max = 5

// This value represents the offset in the stack array. We want to keep this
// number at 2 so that we do not see XErr functions in the stack
const stackFunctionOffset = 2

// XErr - Extended Error struct
type XErr struct {
	Data       map[string]interface{}
	CauseError error
	MaskError  error
	Stack      []StackLocation
}

// XErrEncoded - Extended Error struct for JSON encoding. Used for Marshalling and Unmarshalling
type XErrEncoded struct {
	Data       map[string]interface{} `json:"data"`
	CauseError string                 `json:"causeError"`
	MaskError  string                 `json:"maskError"`
	Stack      []StackLocation        `json:"stack"`
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

// MarshalJSON implements json.Marshaler for XErr
func (xerr *XErr) MarshalJSON() ([]byte, error) {
	x := XErrEncoded{
		Data:       xerr.Data,
		CauseError: xerr.CauseError.Error(),
		MaskError:  xerr.MaskError.Error(),
		Stack:      xerr.Stack,
	}

	return json.Marshal(x)
}

// UnmarshalJSON implements json.Unmarshaler for XErr
func (xerr *XErr) UnmarshalJSON(data []byte) error {
	x := XErrEncoded{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	*xerr = XErr{
		Data:       x.Data,
		CauseError: errors.New(x.CauseError),
		MaskError:  errors.New(x.MaskError),
		Stack:      x.Stack,
	}

	return nil
}

// ToError - Converts XErr into a GoLang error where value is equal to marshalled XErr JSON
func (xerr *XErr) ToError() error {
	if xerr == nil {
		return nil
	}

	if xerr.CauseError == nil {
		return nil
	}

	if xerr.MaskError == nil {
		xerr.MaskError = xerr.CauseError
	}

	result, e := json.Marshal(xerr)
	if e != nil {
		return xerr.MaskError
	}

	return errors.New(string(result))
}

// Extend - Returns a new GoLang error which is a marshalled JSON XErr. It is generated based on the error
// argument passed into the function.
func Extend(err error) error {
	if err == nil {
		return nil
	}

	xerr, ok := ToXErr(err)
	if ok {
		return err
	}

	xerr = &XErr{
		Data:       make(map[string]interface{}),
		CauseError: err,
		MaskError:  err,
		Stack:      getStack(stackFunctionOffset, max),
	}

	return xerr.ToError()
}

// ToXErr - Will return pointer to XErr based on the argument GoLang error if the argument is a marshalled
// XErr error. Second return value is a boolean representing if conversion successful
func ToXErr(err error) (*XErr, bool) {
	if err == nil {
		return nil, false
	}

	xerr := &XErr{}
	if e := json.Unmarshal([]byte(err.Error()), xerr); e != nil {
		return nil, false
	}

	return xerr, true
}

// MaskError - Returns a new GoLang error which is a marshalled JSON XErr. It is generated based on the error
// and mask error arguments passed into the function. If initial error is already a serialized XErr then
// its Mask property will be updated. If Mask is nil then passed in error returned as it is.
func MaskError(err error, mask error) error {
	if err == nil {
		return nil
	}

	if mask == nil {
		return err
	}

	xerr := &XErr{}
	var ok bool

	xerr, ok = ToXErr(err)
	if !ok {
		xerr = &XErr{
			Data:       make(map[string]interface{}),
			CauseError: err,
			MaskError:  mask,
			Stack:      getStack(stackFunctionOffset, max),
		}
	} else {
		xerr.MaskError = mask
	}

	return xerr.ToError()
}

// Cause - Returns Cause Error property of the XErr if the passed error is serialized XErr error
func Cause(err error) error {
	xerr, ok := ToXErr(err)
	if !ok {
		return err
	}

	return xerr.CauseError
}

// GetData - Returns custom Data Error property by name if the passed error is a serialized XErr error
func GetData(err error, name string) (interface{}, bool) {
	xerr, ok := ToXErr(err)
	if !ok {
		return nil, false
	}

	value, ok := xerr.Data[name]
	if !ok {
		return nil, false
	}

	return value, true
}

// SetData - Sets custom Data Error property if the passed error is a serialized XErr error
func SetData(err error, name string, value interface{}) error {
	xerr, ok := ToXErr(err)
	if !ok {
		return err
	}

	xerr.Data[name] = value

	return xerr.ToError()
}

// Error - Is equivalent of calling Error() goroutine of the GoLang error. However if the error
// is a serialized XErr then Mask.Error() is returned back. This function does not expose Cause
// to the caller
func Error(err error) string {
	if err == nil {
		return ""
	}

	xerr, ok := ToXErr(err)
	if !ok {
		return err.Error()
	}

	if xerr.MaskError != nil {
		return xerr.MaskError.Error()
	}

	return ""
}

// IsEqual - Returns true if two supplied errors have the same value. If one of those errors
// XErr then its Cause value will be used for comparing
func IsEqual(err1 error, err2 error) bool {
	var terr1, terr2 error

	if err1 == nil && err2 == nil {
		return true
	}
	if err1 == nil || err2 == nil {
		return false
	}

	xerr1, ok := ToXErr(err1)
	if ok {
		terr1 = xerr1.CauseError
	} else {
		terr1 = err1
	}

	xerr2, ok := ToXErr(err2)
	if ok {
		terr2 = xerr2.CauseError
	} else {
		terr2 = err2
	}

	return terr1.Error() == terr2.Error()
}

// Stack - Returns a printable string of the passed error. If the error is serialized XErr then extra data is added.
// Such data as Cause Error, Mask Error and full Stack if applicable
func Stack(err error) string {
	if err == nil {
		return ""
	}

	xerr, ok := ToXErr(err)
	if !ok {
		return err.Error()
	}

	result := []string{""}

	result = append(result, fmt.Sprintf("[ERROR] %s", xerr.CauseError.Error()))
	if xerr.CauseError != xerr.MaskError {
		result = append(result, fmt.Sprintf("[MASK ERROR] %s", xerr.MaskError.Error()))
	}

	if len(xerr.Stack) == 0 {
		return strings.Join(result, "\n")
	}

	result = append(result, "Stack:")
	for _, stackLocation := range xerr.Stack {
		result = append(result, stackLocation.String())
	}

	return strings.Join(result, "\n")
}

// Returns execution Stack of the goroutine which called it in the form of StackLocation array
// skip - is a starting level on the execution stack where 0 = getStack() function itself, 1 = caller who called getStack(), and so forth
// max - maximum number of stack levels returned
func getStack(skip, max int) []StackLocation {
	stack := []StackLocation{}

	for i := 0; i < max; i++ {
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
	}

	return stack
}
