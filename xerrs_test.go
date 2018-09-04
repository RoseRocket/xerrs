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

// TestExtendError -
func TestExtendError(t *testing.T) {

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
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
		},
		TestCase{
			Description: "Non XErr error",
			Input:       errors.New("ABC"),
			Output:      errors.New("{\"data\":\"\",\"causeError\":\"ABC\",\"maskError\":\"ABC\",\"stack\":[{\"function\":\"xerrs.TestExtendError.func1\",\"file\":\"xerrs_test.go\",\"line\":73},{\"function\":\"convey.parseAction.func1\",\"file\":\"discovery.go\",\"line\":80},{\"function\":\"convey.(*context).conveyInner\",\"file\":\"context.go\",\"line\":261},{\"function\":\"convey.rootConvey.func1\",\"file\":\"context.go\",\"line\":110}]}"),
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := ExtendError(testCase.Input)

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
				Data:       "1",
				CauseError: errors.New("ABC"),
				MaskError:  errors.New("XYZ"),
				Stack:      []StackLocation{},
			},
			Output: errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
		},
		TestCase{
			Description: "Empty Mask",
			Input: &XErr{
				Data:       "1",
				CauseError: errors.New("ABC"),
				Stack:      []StackLocation{},
			},
			Output: errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"ABC\",\"stack\":[]}"),
		},
		TestCase{
			Description: "Non Empty XErr with stack",
			Input: &XErr{
				Data:       "1",
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
			Output: errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[{\"function\":\"blah\",\"file\":\"test.go\",\"line\":1000}]}"),
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
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("ABC"),
		},
		TestCase{
			Description: "wrong format XErr",
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("{\"data\":\"1\",\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := Cause(testCase.Input)
			So(output, ShouldResemble, testCase.Output)
		})
	}
}

// TestData -
func TestData(t *testing.T) {

	type TestCase struct {
		Description string
		Input       error
		OutputData  string
		OutputOK    bool
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			Input:       nil,
			OutputData:  "",
			OutputOK:    false,
		},
		TestCase{
			Description: "Regular error",
			Input:       errors.New("ABC"),
			OutputData:  "",
			OutputOK:    false,
		},
		TestCase{
			Description: "XErr",
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			OutputData:  "1",
			OutputOK:    true,
		},
		TestCase{
			Description: "wrong format XErr",
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			OutputData:  "",
			OutputOK:    false,
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output, ok := Data(testCase.Input)
			So(ok, ShouldEqual, testCase.OutputOK)
			So(output, ShouldResemble, testCase.OutputData)
		})
	}
}

// TestSetData -
func TestSetData(t *testing.T) {

	type TestCase struct {
		Description string
		InputData   string
		InputError  error
		Output      error
	}

	testCases := []TestCase{
		TestCase{
			Description: "nil error",
			InputData:   "SOME_DATA",
			InputError:  nil,
			Output:      nil,
		},
		TestCase{
			Description: "Regular error",
			InputData:   "SOME_DATA",
			InputError:  errors.New("ABC"),
			Output:      errors.New("ABC"),
		},
		TestCase{
			Description: "XErr",
			InputData:   "SOME_DATA",
			InputError:  errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("{\"data\":\"SOME_DATA\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
		},
		TestCase{
			Description: "wrong format XErr",
			InputData:   "SOME_DATA",
			InputError:  errors.New("{\"data\":\"1\",\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			Output:      errors.New("{\"data\":\"1\",\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := SetData(testCase.InputError, testCase.InputData)
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
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			Output:      "XYZ",
		},
		TestCase{
			Description: "wrong format XErr",
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
			Output:      "{\"data\":\"1\",\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}",
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
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
			OutputXErr: &XErr{
				Data:       "1",
				CauseError: errors.New("ABC"),
				MaskError:  errors.New("XYZ"),
				Stack:      []StackLocation{},
			},
			OutputOK: true,
		},
		TestCase{
			Description: "wrong format XErr",
			Input:       errors.New("{\"data\":\"1\",\"causeError\":\"ABC\"\"maskErr:\"XYZ\",\"stack\":[]}"),
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
			Output:      errors.New("{\"data\":\"\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[{\"function\":\"xerrs.TestMask.func1\",\"file\":\"xerrs_test.go\",\"line\":408},{\"function\":\"convey.parseAction.func1\",\"file\":\"discovery.go\",\"line\":80},{\"function\":\"convey.(*context).conveyInner\",\"file\":\"context.go\",\"line\":261},{\"function\":\"convey.rootConvey.func1\",\"file\":\"context.go\",\"line\":110}]}"),
		},
		TestCase{
			Description: "XErr",
			InputErr:    errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"BLAH\",\"stack\":[]}"),
			InputMask:   errors.New("XYZ"),
			Output:      errors.New("{\"data\":\"1\",\"causeError\":\"ABC\",\"maskError\":\"XYZ\",\"stack\":[]}"),
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
					Line:     252,
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
					Line:     501,
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
					Line:     501,
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
