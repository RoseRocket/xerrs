package xerrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const max = 4

// XErr - Extended Error struct
type XErr struct {
	Data       string
	CauseError error
	MaskError  error
	Stack      []StackLocation
}

// XErrEncoded - Extended Error struct for JSON encoding. Used for Marshalling and Unmarshalling
type XErrEncoded struct {
	Data       string          `json:"data"`
	CauseError string          `json:"causeError"`
	MaskError  string          `json:"maskError"`
	Stack      []StackLocation `json:"stack"`
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
	x := XErrEncoded{}

	x.Data = xerr.Data
	x.CauseError = xerr.CauseError.Error()
	x.MaskError = xerr.MaskError.Error()
	x.Stack = xerr.Stack

	return json.Marshal(x)
}

// UnmarshalJSON implements json.Unmarshaler for XErr
func (xerr *XErr) UnmarshalJSON(data []byte) error {
	x := XErrEncoded{}

	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}

	*xerr = XErr{}

	xerr.Data = x.Data
	xerr.CauseError = errors.New(x.CauseError)
	xerr.MaskError = errors.New(x.MaskError)
	xerr.Stack = x.Stack

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

// ExtendError - Returns a new GoLang error which is a marshalled JSON XErr. It is generated based on the error
// argument passed into the function.
func ExtendError(err error) error {
	if err == nil {
		return nil
	}

	xerr, ok := ToXErr(err)
	if ok {
		return err
	}

	xerr = &XErr{
		Data:       "",
		CauseError: err,
		MaskError:  err,
		Stack:      getStack(2, max),
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
	e := json.Unmarshal([]byte(err.Error()), xerr)
	if e != nil {
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
			Data:       "",
			CauseError: err,
			MaskError:  mask,
			Stack:      getStack(2, max),
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

// Data - Returns Data Error property of the XErr if the passed error is serialized XErr error
func Data(err error) (string, bool) {
	xerr, ok := ToXErr(err)
	if !ok {
		return "", false
	}

	return xerr.Data, true
}

// SetData - Sets Data Error property of the XErr if the passed error is serialized XErr error
func SetData(err error, data string) error {
	xerr, ok := ToXErr(err)
	if !ok {
		return err
	}

	xerr.Data = data

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
