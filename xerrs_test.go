package xerrs

import (
	"errors"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// get the last part of the path
func getLastPathPart(str string) string {
	parts := strings.Split(str, "/")

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

// Remove the path to the files in stack so that Testing would not be path dependent
func transformStack(stack []StackLocation) []StackLocation {
	for i := range stack {
		stack[i].Function = getLastPathPart(stack[i].Function)
		stack[i].File = getLastPathPart(stack[i].File)
	}

	return stack
}

// Transfor XErr by changing Stack not to be path dependent
func transformXErr(err error) error {
	xerr, ok := ToXErr(err)
	if !ok {
		return err
	}

	xerr.Stack = transformStack(xerr.Stack)

	return xerr.ToError()
}

// TestExtend -
func TestExtend(t *testing.T) {

	type TestCase struct {
		Description string
		Input       error
		Output      error
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			Input:       nil,
			Output:      nil,
		},
		TestCase{
			Description: "XErr error",
			Input:       errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
		},
		TestCase{
			Description: "Non XErr error",
			Input:       errors.New("ABC"),
			Output:      errors.New("{\"data\":{},\"causeError\":\"ABC\",\"maskError\":\"ABC\",\"stack\":[{\"function\":\"xerrs.TestExtend.func1\",\"file\":\"xerrs_test.go\",\"line\":73},{\"function\":\"convey.parseAction.func1\",\"file\":\"discovery.go\",\"line\":80},{\"function\":\"convey.(*context).conveyInner\",\"file\":\"context.go\",\"line\":261},{\"function\":\"convey.rootConvey.func1\",\"file\":\"context.go\",\"line\":110},{\"function\":\"gls.(*ContextManager).SetValues.func1\",\"file\":\"context.go\",\"line\":97}]}"),
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := Extend(testCase.Input)

			output = transformXErr(output)

			So(output, ShouldResemble, testCase.Output)
		})
	}
}

// TestToError -
func TestToError(t *testing.T) {

	type TestCase struct {
		Description string
		Input       *XErr
		Output      error
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			Input:       nil,
			Output:      nil,
		},
		TestCase{
			Description: "Cause is nil",
			Input:       &XErr{},
			Output:      nil,
		},
		TestCase{
			Description: "Non Empty XErr",
			Input: &XErr{
				Data: map[string]interface{}{
					"SOME_DATA": 1,
				},
				CauseError: errors.New("ABC"),
				MaskError:  errors.New("XYZ"),
				Stack:      []StackLocation{},
			},
			Output: errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
		},
		TestCase{
			Description: "Empty Mask",
			Input: &XErr{
				Data: map[string]interface{}{
					"SOME_DATA": 1,
				},
				CauseError: errors.New("ABC"),
				Stack:      []StackLocation{},
			},
			Output: errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"ABC\",\"stack\":[]}"),
		},
		TestCase{
			Description: "Non Empty XErr with stack",
			Input: &XErr{
				Data: map[string]interface{}{
					"SOME_DATA": 1,
				},
				CauseError: errors.New("ABC"),
				MaskError:  errors.New("XYZ"),
				Stack: []StackLocation{
					StackLocation{
						Function: "blah",
						File:     "test.go",
						Line:     1000,
					},
				},
			},
			Output: errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[{\"function\":\"blah\",\"file\":\"test.go\",\"line\":1000}]}"),
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := testCase.Input.ToError()
			So(output, ShouldResemble, testCase.Output)
		})
	}
}

// TestCause -
func TestCause(t *testing.T) {

	type TestCase struct {
		Description string
		Input       error
		Output      error
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			Input:       nil,
			Output:      nil,
		},
		TestCase{
			Description: "Regular error",
			Input:       errors.New("ABC"),
			Output:      errors.New("ABC"),
		},
		TestCase{
			Description: "XErr",
			Input:       errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("ABC"),
		},
		TestCase{
			Description: "wrong format XErr",
			Input:       errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := Cause(testCase.Input)
			So(output, ShouldResemble, testCase.Output)
		})
	}
}

// TestGetData -
func TestGetData(t *testing.T) {

	type TestCase struct {
		Description string
		InputErr    error
		InputName   string
		OutputData  interface{}
		OutputOK    bool
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			InputErr:    nil,
			InputName:   "SOME_DATA",
			OutputData:  nil,
			OutputOK:    false,
		},
		TestCase{
			Description: "Regular error",
			InputErr:    errors.New("ABC"),
			InputName:   "SOME_DATA",
			OutputData:  nil,
			OutputOK:    false,
		},
		TestCase{
			Description: "XErr existing key",
			InputErr:    errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			InputName:   "SOME_DATA",
			OutputData:  1,
			OutputOK:    true,
		},
		TestCase{
			Description: "XErr non existing key",
			InputErr:    errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			InputName:   "OTHER_KEY",
			OutputData:  nil,
			OutputOK:    false,
		},
		TestCase{
			Description: "wrong format XErr",
			InputErr:    errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			InputName:   "SOME_DATA",
			OutputData:  nil,
			OutputOK:    false,
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output, ok := GetData(testCase.InputErr, testCase.InputName)
			So(ok, ShouldEqual, testCase.OutputOK)
			So(output, ShouldEqual, testCase.OutputData)
		})
	}
}

// TestSetData -
func TestSetData(t *testing.T) {

	type TestCase struct {
		Description string
		InputName   string
		InputValue  interface{}
		InputError  error
		Output      error
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			InputName:   "SOME_DATA",
			InputValue:  1,
			InputError:  nil,
			Output:      nil,
		},
		TestCase{
			Description: "Regular error",
			InputName:   "SOME_DATA",
			InputValue:  1,
			InputError:  errors.New("ABC"),
			Output:      errors.New("ABC"),
		},
		TestCase{
			Description: "XErr",
			InputName:   "SOME_DATA",
			InputValue:  1,
			InputError:  errors.New("{\"data\":{},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
		},
		TestCase{
			Description: "wrong format XErr",
			InputName:   "SOME_DATA",
			InputValue:  1,
			InputError:  errors.New("{\"data\":{},\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("{\"data\":{},\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := SetData(testCase.InputError, testCase.InputName, testCase.InputValue)
			So(output, ShouldResemble, testCase.Output)
		})
	}
}

// TestError -
func TestError(t *testing.T) {

	type TestCase struct {
		Description string
		Input       error
		Output      string
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			Input:       nil,
			Output:      "",
		},
		TestCase{
			Description: "Regular error",
			Input:       errors.New("ABC"),
			Output:      "ABC",
		},
		TestCase{
			Description: "XErr",
			Input:       errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			Output:      "XYZ",
		},
		TestCase{
			Description: "wrong format XErr",
			Input:       errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			Output:      "{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}",
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := Error(testCase.Input)
			So(output, ShouldResemble, testCase.Output)
		})
	}
}

// TestToXErr -
func TestToXErr(t *testing.T) {

	type TestCase struct {
		Description string
		Input       error
		OutputXErr  *XErr
		OutputOK    bool
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			Input:       nil,
			OutputXErr:  nil,
			OutputOK:    false,
		},
		TestCase{
			Description: "Regular error",
			Input:       errors.New("ABC"),
			OutputXErr:  nil,
			OutputOK:    false,
		},
		TestCase{
			Description: "XErr",
			Input:       errors.New("{\"data\":{},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			OutputXErr: &XErr{
				Data:       map[string]interface{}{},
				CauseError: errors.New("ABC"),
				MaskError:  errors.New("XYZ"),
				Stack:      []StackLocation{},
			},
			OutputOK: true,
		},
		TestCase{
			Description: "wrong format XErr",
			Input:       errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			OutputXErr:  nil,
			OutputOK:    false,
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output, ok := ToXErr(testCase.Input)
			So(ok, ShouldEqual, testCase.OutputOK)
			So(output, ShouldResemble, testCase.OutputXErr)
		})
	}
}

// TestMask -
func TestMask(t *testing.T) {

	type TestCase struct {
		Description string
		InputErr    error
		InputMask   error
		Output      error
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			InputErr:    nil,
			InputMask:   nil,
			Output:      nil,
		},
		TestCase{
			Description: "mask is nil",
			InputErr:    errors.New("ABC"),
			InputMask:   nil,
			Output:      errors.New("ABC"),
		},
		TestCase{
			Description: "Regular error",
			InputErr:    errors.New("ABC"),
			InputMask:   errors.New("XYZ"),
			Output:      errors.New("{\"data\":{},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[{\"function\":\"xerrs.TestMask.func1\",\"file\":\"xerrs_test.go\",\"line\":431},{\"function\":\"convey.parseAction.func1\",\"file\":\"discovery.go\",\"line\":80},{\"function\":\"convey.(*context).conveyInner\",\"file\":\"context.go\",\"line\":261},{\"function\":\"convey.rootConvey.func1\",\"file\":\"context.go\",\"line\":110},{\"function\":\"gls.(*ContextManager).SetValues.func1\",\"file\":\"context.go\",\"line\":97}]}"),
		},
		TestCase{
			Description: "XErr",
			InputErr:    errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			InputMask:   errors.New("XYZ"),
			Output:      errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := MaskError(testCase.InputErr, testCase.InputMask)

			output = transformXErr(output)

			So(output, ShouldResemble, testCase.Output)
		})
	}
}

// TestGetStack -
func TestGetStack(t *testing.T) {

	type TestCase struct {
		Description string
		InputSkip   int
		InputMax    int
		Output      []StackLocation
	}

	testCases := []TestCase{
		TestCase{
			Description: "0 skip and 0 max",
			InputSkip:   0,
			InputMax:    0,
			Output:      []StackLocation{},
		},
		TestCase{
			Description: "1 skip and 0 max",
			InputSkip:   1,
			InputMax:    0,
			Output:      []StackLocation{},
		},
		TestCase{
			Description: "0 skip and 1 max",
			InputSkip:   0,
			InputMax:    1,
			Output: []StackLocation{
				StackLocation{
					Function: "xerrs.getStack",
					File:     "xerrs.go",
					Line:     289,
				},
			},
		},
		TestCase{
			Description: "1 skip and 1 max",
			InputSkip:   1,
			InputMax:    1,
			Output: []StackLocation{
				StackLocation{
					Function: "xerrs.TestGetStack.func1",
					File:     "xerrs_test.go",
					Line:     524,
				},
			},
		},
		TestCase{
			Description: "1 skip and 4 max",
			InputSkip:   1,
			InputMax:    4,
			Output: []StackLocation{
				StackLocation{
					Function: "xerrs.TestGetStack.func1",
					File:     "xerrs_test.go",
					Line:     524,
				},
				StackLocation{
					Function: "convey.parseAction.func1",
					File:     "discovery.go",
					Line:     80,
				},
				StackLocation{
					Function: "convey.(*context).conveyInner",
					File:     "context.go",
					Line:     261,
				},
				StackLocation{
					Function: "convey.rootConvey.func1",
					File:     "context.go",
					Line:     110,
				},
			},
		},
		TestCase{
			Description: "80 skip and 4 max",
			InputSkip:   80,
			InputMax:    4,
			Output:      []StackLocation{},
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := getStack(testCase.InputSkip, testCase.InputMax)

			output = transformStack(output)

			So(output, ShouldResemble, testCase.Output)
		})
	}
}

// TestIsEqual -
func TestIsEqual(t *testing.T) {
	type TestCase struct {
		Description string
		InputErr1   error
		InputErr2   error
		Output      bool
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			InputErr1:   nil,
			InputErr2:   nil,
			Output:      true,
		},
		TestCase{
			Description: "one error is nil",
			InputErr1:   errors.New("ABC"),
			InputErr2:   nil,
			Output:      false,
		},
		TestCase{
			Description: "one error is nil",
			InputErr1:   nil,
			InputErr2:   errors.New("ABC"),
			Output:      false,
		},
		TestCase{
			Description: "Regular errors",
			InputErr1:   errors.New("ABC"),
			InputErr2:   errors.New("XYZ"),
			Output:      false,
		},
		TestCase{
			Description: "one is XErr both equal",
			InputErr1:   errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			InputErr2:   errors.New("ABC"),
			Output:      true,
		},
		TestCase{
			Description: "another is XErr both equal",
			InputErr1:   errors.New("ABC"),
			InputErr2:   errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			Output:      true,
		},
		TestCase{
			Description: "both are XErr both equal",
			InputErr1:   errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			InputErr2:   errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			Output:      true,
		},
		TestCase{
			Description: "one is XErr both not equal",
			InputErr1:   errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			InputErr2:   errors.New("XYZ"),
			Output:      false,
		},
		TestCase{
			Description: "another is XErr both not equal",
			InputErr1:   errors.New("XYZ"),
			InputErr2:   errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			Output:      false,
		},
		TestCase{
			Description: "both are XErr both not equal",
			InputErr1:   errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			InputErr2:   errors.New("{\"data\":{\"SOME_DATA\":1},\"causeError\":\"XYZ\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			Output:      false,
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := IsEqual(testCase.InputErr1, testCase.InputErr2)

			So(output, ShouldEqual, testCase.Output)
		})
	}
}
