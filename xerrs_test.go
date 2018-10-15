package xerrs

import (
	"errors"
	"fmt"
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

// TestNew -
func TestNew(t *testing.T) {
	type TestCase struct {
		Description string
		Input       string
		Output      []StackLocation
	}

	testCases := []TestCase{
		TestCase{
			Description: "some error",
			Input:       "some error",
			Output: []StackLocation{
				StackLocation{
					Function: "xerrs.TestNew.func1",
					File:     "xerrs_test.go",
					Line:     132,
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
				StackLocation{
					Function: "gls.(*ContextManager).SetValues.func1",
					File:     "context.go",
					Line:     97,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId.func1",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls._m",
					File:     "stack_tags.go",
					Line:     74,
				},
				StackLocation{
					Function: "gls.github_com_jtolds_gls_markS",
					File:     "stack_tags.go",
					Line:     54,
				},
				StackLocation{
					Function: "gls.addStackTag",
					File:     "stack_tags.go",
					Line:     49,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls.(*ContextManager).SetValues",
					File:     "context.go",
					Line:     63,
				},
				StackLocation{
					Function: "convey.rootConvey",
					File:     "context.go",
					Line:     105,
				},
				StackLocation{
					Function: "convey.Convey",
					File:     "doc.go",
					Line:     75,
				},
				StackLocation{
					Function: "xerrs.TestNew",
					File:     "xerrs_test.go",
					Line:     131,
				},
				StackLocation{
					Function: "testing.tRunner",
					File:     "testing.go",
					Line:     827,
				},
				StackLocation{
					Function: "runtime.goexit",
					File:     "asm_amd64.s",
					Line:     1333,
				},
			},
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := New(testCase.Input)

			if x, ok := output.(*xerr); ok {
				So(x.Cause().Error(), ShouldEqual, testCase.Input)

				stack := x.Stack()
				stack = transformStack(stack)

				So(stack, ShouldResemble, testCase.Output)
			}

			So(output.Error(), ShouldEqual, testCase.Input)
		})
	}
}

// TestErrorf -
func TestErrorf(t *testing.T) {
	type TestCase struct {
		Description string
		Input       string
		InputArgs   []interface{}
		Output      []StackLocation
	}

	testCases := []TestCase{
		TestCase{
			Description: "some error",
			Input:       "some error %d %v",
			InputArgs: []interface{}{
				1,
				"Hello",
			},
			Output: []StackLocation{
				StackLocation{
					Function: "xerrs.TestErrorf.func1",
					File:     "xerrs_test.go",
					Line:     254,
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
				StackLocation{
					Function: "gls.(*ContextManager).SetValues.func1",
					File:     "context.go",
					Line:     97,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId.func1",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls._m",
					File:     "stack_tags.go",
					Line:     74,
				},
				StackLocation{
					Function: "gls.github_com_jtolds_gls_markS",
					File:     "stack_tags.go",
					Line:     54,
				},
				StackLocation{
					Function: "gls.addStackTag",
					File:     "stack_tags.go",
					Line:     49,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls.(*ContextManager).SetValues",
					File:     "context.go",
					Line:     63,
				},
				StackLocation{
					Function: "convey.rootConvey",
					File:     "context.go",
					Line:     105,
				},
				StackLocation{
					Function: "convey.Convey",
					File:     "doc.go",
					Line:     75,
				},
				StackLocation{
					Function: "xerrs.TestErrorf",
					File:     "xerrs_test.go",
					Line:     251,
				},
				StackLocation{
					Function: "testing.tRunner",
					File:     "testing.go",
					Line:     827,
				},
				StackLocation{
					Function: "runtime.goexit",
					File:     "asm_amd64.s",
					Line:     1333,
				},
			},
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			message := fmt.Sprintf(testCase.Input, testCase.InputArgs...)

			output := Errorf(testCase.Input, testCase.InputArgs...)

			if x, ok := output.(*xerr); ok {
				So(x.Cause().Error(), ShouldEqual, message)

				stack := x.Stack()
				stack = transformStack(stack)

				So(stack, ShouldResemble, testCase.Output)
			}

			So(output.Error(), ShouldEqual, message)
		})
	}
}

// TestExtend -
func TestExtend(t *testing.T) {
	type TestCase struct {
		Description string
		Input       error
		Output      []StackLocation
	}

	testCases := []TestCase{
		TestCase{
			Description: "some error",
			Input:       errors.New("ERROR HERE"),
			Output: []StackLocation{
				StackLocation{
					Function: "xerrs.TestExtend.func1",
					File:     "xerrs_test.go",
					Line:     374,
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
				StackLocation{
					Function: "gls.(*ContextManager).SetValues.func1",
					File:     "context.go",
					Line:     97,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId.func1",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls._m",
					File:     "stack_tags.go",
					Line:     74,
				},
				StackLocation{
					Function: "gls.github_com_jtolds_gls_markS",
					File:     "stack_tags.go",
					Line:     54,
				},
				StackLocation{
					Function: "gls.addStackTag",
					File:     "stack_tags.go",
					Line:     49,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls.(*ContextManager).SetValues",
					File:     "context.go",
					Line:     63,
				},
				StackLocation{
					Function: "convey.rootConvey",
					File:     "context.go",
					Line:     105,
				},
				StackLocation{
					Function: "convey.Convey",
					File:     "doc.go",
					Line:     75,
				},
				StackLocation{
					Function: "xerrs.TestExtend",
					File:     "xerrs_test.go",
					Line:     373,
				},
				StackLocation{
					Function: "testing.tRunner",
					File:     "testing.go",
					Line:     827,
				},
				StackLocation{
					Function: "runtime.goexit",
					File:     "asm_amd64.s",
					Line:     1333,
				},
			},
		},
		TestCase{
			Description: "nil error",
			Input:       nil,
			Output:      []StackLocation{},
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := Extend(testCase.Input)

			if testCase.Input != nil {
				if x, ok := output.(*xerr); ok {
					So(x.Cause().Error(), ShouldEqual, testCase.Input.Error())

					stack := x.Stack()
					stack = transformStack(stack)

					So(stack, ShouldResemble, testCase.Output)
				}

				So(output.Error(), ShouldEqual, testCase.Input.Error())
			} else {
				So(output, ShouldEqual, nil)
			}

		})
	}
}

// TestMask -
func TestMask(t *testing.T) {
	type TestCase struct {
		Description string
		InputErr    error
		MaskErr     error
		Output      []StackLocation
	}

	testCases := []TestCase{
		TestCase{
			Description: "some error",
			InputErr:    errors.New("ERROR HERE"),
			MaskErr:     errors.New("MASK ERROR HERE"),
			Output: []StackLocation{
				StackLocation{
					Function: "xerrs.TestMask.func1",
					File:     "xerrs_test.go",
					Line:     589,
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
				StackLocation{
					Function: "gls.(*ContextManager).SetValues.func1",
					File:     "context.go",
					Line:     97,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId.func1",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls._m",
					File:     "stack_tags.go",
					Line:     74,
				},
				StackLocation{
					Function: "gls.github_com_jtolds_gls_markS",
					File:     "stack_tags.go",
					Line:     54,
				},
				StackLocation{
					Function: "gls.addStackTag",
					File:     "stack_tags.go",
					Line:     49,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls.(*ContextManager).SetValues",
					File:     "context.go",
					Line:     63,
				},
				StackLocation{
					Function: "convey.rootConvey",
					File:     "context.go",
					Line:     105,
				},
				StackLocation{
					Function: "convey.Convey",
					File:     "doc.go",
					Line:     75,
				},
				StackLocation{
					Function: "xerrs.TestMask",
					File:     "xerrs_test.go",
					Line:     588,
				},
				StackLocation{
					Function: "testing.tRunner",
					File:     "testing.go",
					Line:     827,
				},
				StackLocation{
					Function: "runtime.goexit",
					File:     "asm_amd64.s",
					Line:     1333,
				},
			},
		},
		TestCase{
			Description: "nil error",
			InputErr:    nil,
			MaskErr:     nil,
			Output:      []StackLocation{},
		},
		TestCase{
			Description: "nil mask",
			InputErr:    errors.New("ERROR HERE"),
			MaskErr:     nil,
			Output: []StackLocation{
				StackLocation{
					Function: "xerrs.TestMask.func1",
					File:     "xerrs_test.go",
					Line:     589,
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
				StackLocation{
					Function: "gls.(*ContextManager).SetValues.func1",
					File:     "context.go",
					Line:     97,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId.func1",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls._m",
					File:     "stack_tags.go",
					Line:     74,
				},
				StackLocation{
					Function: "gls.github_com_jtolds_gls_markS",
					File:     "stack_tags.go",
					Line:     54,
				},
				StackLocation{
					Function: "gls.addStackTag",
					File:     "stack_tags.go",
					Line:     49,
				},
				StackLocation{
					Function: "gls.EnsureGoroutineId",
					File:     "gid.go",
					Line:     24,
				},
				StackLocation{
					Function: "gls.(*ContextManager).SetValues",
					File:     "context.go",
					Line:     63,
				},
				StackLocation{
					Function: "convey.rootConvey",
					File:     "context.go",
					Line:     105,
				},
				StackLocation{
					Function: "convey.Convey",
					File:     "doc.go",
					Line:     75,
				},
				StackLocation{
					Function: "xerrs.TestMask",
					File:     "xerrs_test.go",
					Line:     588,
				},
				StackLocation{
					Function: "testing.tRunner",
					File:     "testing.go",
					Line:     827,
				},
				StackLocation{
					Function: "runtime.goexit",
					File:     "asm_amd64.s",
					Line:     1333,
				},
			},
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := Mask(testCase.InputErr, testCase.MaskErr)

			if testCase.InputErr != nil {
				if x, ok := output.(*xerr); ok {
					So(x.Cause().Error(), ShouldEqual, testCase.InputErr.Error())

					stack := x.Stack()
					stack = transformStack(stack)

					So(stack, ShouldResemble, testCase.Output)
				}

				if testCase.MaskErr != nil {
					So(output.Error(), ShouldNotEqual, testCase.InputErr.Error())
					So(output.Error(), ShouldEqual, testCase.MaskErr.Error())
				} else {
					So(output.Error(), ShouldEqual, testCase.InputErr.Error())
				}

			} else {
				So(output, ShouldEqual, nil)
			}

		})
	}
}

// TestXMask -
func TestXMask(t *testing.T) {
	type TestCase struct {
		Description string
		InputStr    string
		InputMask   error
		Output      string
	}

	testCases := []TestCase{
		TestCase{
			Description: "basic error",
			InputStr:    "ERROR",
			InputMask:   errors.New("MASK ERROR"),
			Output:      "MASK ERROR",
		},
		TestCase{
			Description: "nil mask",
			InputStr:    "ERROR",
			InputMask:   nil,
			Output:      "ERROR",
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := New(testCase.InputStr)

			if x, ok := output.(*xerr); ok {
				x.Mask(testCase.InputMask)
			}

			So(output.Error(), ShouldEqual, testCase.Output)
		})
	}
}

// TestGetDataAndSetData -
func TestGetDataAndSetData(t *testing.T) {
	Convey("GetData() and SetData()", t, func() {
		output := New("ERROR")

		if x, ok := output.(*xerr); ok {
			val, ok := x.GetData("SOME_DATA")
			So(ok, ShouldEqual, false)

			x.SetData("SOME_DATA", 100)

			val, ok = x.GetData("SOME_DATA")
			So(ok, ShouldEqual, true)
			So(val, ShouldEqual, 100)
		}

		So(output.Error(), ShouldEqual, "ERROR")
	})
}

// TestDetails -
func TestDetails(t *testing.T) {
	type TestCase struct {
		Description   string
		InputError    error
		InputMask     error
		InputMaxStack int
		Output        string
	}

	testCases := []TestCase{
		TestCase{
			Description:   "basic",
			InputError:    errors.New("ERROR"),
			InputMask:     errors.New("MASK"),
			InputMaxStack: 100,
			Output: `
[ERROR] ERROR
[MASK ERROR] MASK
[STACK]:
xerrs.TestDetails.func1 [xerrs_test.go:778]
convey.parseAction.func1 [discovery.go:80]
convey.(*context).conveyInner [context.go:261]
convey.rootConvey.func1 [context.go:110]
gls.(*ContextManager).SetValues.func1 [context.go:97]
gls.EnsureGoroutineId.func1 [gid.go:24]
gls._m [stack_tags.go:74]
gls.github_com_jtolds_gls_markS [stack_tags.go:54]
gls.addStackTag [stack_tags.go:49]
gls.EnsureGoroutineId [gid.go:24]
gls.(*ContextManager).SetValues [context.go:63]
convey.rootConvey [context.go:105]
convey.Convey [doc.go:75]
xerrs.TestDetails [xerrs_test.go:777]
testing.tRunner [testing.go:827]
runtime.goexit [asm_amd64.s:1333]`,
		},
		TestCase{
			Description:   "mask is the same as error",
			InputError:    errors.New("ERROR"),
			InputMask:     errors.New("ERROR"),
			InputMaxStack: 100,
			Output: `
[ERROR] ERROR
[STACK]:
xerrs.TestDetails.func1 [xerrs_test.go:778]
convey.parseAction.func1 [discovery.go:80]
convey.(*context).conveyInner [context.go:261]
convey.rootConvey.func1 [context.go:110]
gls.(*ContextManager).SetValues.func1 [context.go:97]
gls.EnsureGoroutineId.func1 [gid.go:24]
gls._m [stack_tags.go:74]
gls.github_com_jtolds_gls_markS [stack_tags.go:54]
gls.addStackTag [stack_tags.go:49]
gls.EnsureGoroutineId [gid.go:24]
gls.(*ContextManager).SetValues [context.go:63]
convey.rootConvey [context.go:105]
convey.Convey [doc.go:75]
xerrs.TestDetails [xerrs_test.go:777]
testing.tRunner [testing.go:827]
runtime.goexit [asm_amd64.s:1333]`,
		},
		TestCase{
			Description:   "mask is nil",
			InputError:    errors.New("ERROR"),
			InputMask:     nil,
			InputMaxStack: 100,
			Output: `
[ERROR] ERROR
[STACK]:
xerrs.TestDetails.func1 [xerrs_test.go:778]
convey.parseAction.func1 [discovery.go:80]
convey.(*context).conveyInner [context.go:261]
convey.rootConvey.func1 [context.go:110]
gls.(*ContextManager).SetValues.func1 [context.go:97]
gls.EnsureGoroutineId.func1 [gid.go:24]
gls._m [stack_tags.go:74]
gls.github_com_jtolds_gls_markS [stack_tags.go:54]
gls.addStackTag [stack_tags.go:49]
gls.EnsureGoroutineId [gid.go:24]
gls.(*ContextManager).SetValues [context.go:63]
convey.rootConvey [context.go:105]
convey.Convey [doc.go:75]
xerrs.TestDetails [xerrs_test.go:777]
testing.tRunner [testing.go:827]
runtime.goexit [asm_amd64.s:1333]`,
		},
		TestCase{
			Description:   "fewer stack lines",
			InputError:    errors.New("ERROR"),
			InputMask:     errors.New("MASK"),
			InputMaxStack: 4,
			Output: `
[ERROR] ERROR
[MASK ERROR] MASK
[STACK]:
xerrs.TestDetails.func1 [xerrs_test.go:778]
convey.parseAction.func1 [discovery.go:80]
convey.(*context).conveyInner [context.go:261]
convey.rootConvey.func1 [context.go:110]`,
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			output := Extend(testCase.InputError)

			if x, ok := output.(*xerr); ok {
				x.Mask(testCase.InputMask)

				stack := x.Stack()
				stack = transformStack(stack)

				So(x.Details(testCase.InputMaxStack), ShouldEqual, testCase.Output)
			}
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
			Description: "both errors are nil",
			InputErr1:   nil,
			InputErr2:   nil,
			Output:      true,
		},
		TestCase{
			Description: "one errors is nil",
			InputErr1:   errors.New("ABC"),
			InputErr2:   nil,
			Output:      false,
		},
		TestCase{
			Description: "both errors are nil",
			InputErr1:   nil,
			InputErr2:   errors.New("ABC"),
			Output:      false,
		},
		TestCase{
			Description: "both errors are basic ones. equal",
			InputErr1:   errors.New("ABC"),
			InputErr2:   errors.New("ABC"),
			Output:      true,
		},
		TestCase{
			Description: "both errors are basic ones. not equal",
			InputErr1:   errors.New("ABC"),
			InputErr2:   errors.New("XYZ"),
			Output:      false,
		},
		TestCase{
			Description: "one errors is xerr another is basic ones. equal",
			InputErr1:   Extend(errors.New("ABC")),
			InputErr2:   errors.New("ABC"),
			Output:      true,
		},
		TestCase{
			Description: "one errors is xerr another is basic ones. not equal",
			InputErr1:   Extend(errors.New("XYZ")),
			InputErr2:   errors.New("ABC"),
			Output:      false,
		},
		TestCase{
			Description: "one errors is xerr another is basic ones. different mask. equal",
			InputErr1:   Mask(errors.New("ABC"), errors.New("XYZ")),
			InputErr2:   errors.New("ABC"),
			Output:      true,
		},
		TestCase{
			Description: "both errors are xerr. equal",
			InputErr1:   Extend(errors.New("ABC")),
			InputErr2:   Extend(errors.New("ABC")),
			Output:      true,
		},
		TestCase{
			Description: "both errors are xerr. not equal",
			InputErr1:   Extend(errors.New("XYZ")),
			InputErr2:   Extend(errors.New("ABC")),
			Output:      false,
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			So(IsEqual(testCase.InputErr1, testCase.InputErr2), ShouldEqual, testCase.Output)
		})
	}
}
