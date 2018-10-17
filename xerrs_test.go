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

// TestNew -
func TestNew(t *testing.T) {
	Convey("Testing New function", t, func() {
		output := New("ABC")

		if x, ok := output.(*xerr); ok {
			So(x.cause.Error(), ShouldEqual, "ABC")
			So(x.mask, ShouldEqual, nil)

			stack := transformStack(x.stack)

			So(stack, ShouldResemble, []StackLocation{
				StackLocation{
					Function: "xerrs.TestNew.func1",
					File:     "xerrs_test.go",
					Line:     35,
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
					Line:     34,
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
			})
		}

		So(output.Error(), ShouldEqual, "ABC")
	})
}

// TestErrorf -
func TestErrorf(t *testing.T) {
	Convey("Testing Errorf function", t, func() {
		output := Errorf("some error %d %v", 1, "HELLO")

		if x, ok := output.(*xerr); ok {
			So(x.cause.Error(), ShouldEqual, "some error 1 HELLO")
			So(x.mask, ShouldEqual, nil)

			stack := transformStack(x.stack)

			So(stack, ShouldResemble, []StackLocation{
				StackLocation{
					Function: "xerrs.TestErrorf.func1",
					File:     "xerrs_test.go",
					Line:     134,
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
					Line:     133,
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
			})
		}

		So(output.Error(), ShouldEqual, "some error 1 HELLO")
	})
}

// TestExtend -
func TestExtend(t *testing.T) {
	Convey("Testing Extend function", t, func() {
		output := Extend(errors.New("ABC"))

		if x, ok := output.(*xerr); ok {
			So(x.cause.Error(), ShouldEqual, "ABC")
			So(x.mask, ShouldEqual, nil)

			stack := transformStack(x.stack)

			So(stack, ShouldResemble, []StackLocation{
				StackLocation{
					Function: "xerrs.TestExtend.func1",
					File:     "xerrs_test.go",
					Line:     233,
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
					Line:     232,
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
			})
		}

		So(output.Error(), ShouldEqual, "ABC")
	})
}

// TestMask -
func TestMask(t *testing.T) {
	Convey("nil error", t, func() {
		err := Mask(nil, errors.New("ABC"))
		So(err, ShouldEqual, nil)
	})

	Convey("basic error with nil mask", t, func() {
		err := Mask(errors.New("ABC"), nil)
		_, ok := err.(*xerr)
		So(err.Error(), ShouldEqual, "ABC")
		So(ok, ShouldEqual, true)
	})

	Convey("basic error with not nil mask", t, func() {
		err := Mask(errors.New("ABC"), errors.New("XYZ"))
		_, ok := err.(*xerr)
		So(err.Error(), ShouldEqual, "XYZ")
		So(ok, ShouldEqual, true)
	})

	Convey("xerr without a mask", t, func() {
		intial := Extend(errors.New("ABC"))
		err := Mask(intial, errors.New("XYZ"))
		_, ok := err.(*xerr)
		So(err.Error(), ShouldEqual, "XYZ")
		So(ok, ShouldEqual, true)
	})

	Convey("xerr with a mask", t, func() {
		intial := Mask(errors.New("ABC"), errors.New("001"))
		err := Mask(intial, errors.New("XYZ"))
		_, ok := err.(*xerr)
		So(err.Error(), ShouldEqual, "XYZ")
		So(ok, ShouldEqual, true)
	})

	Convey("xerr setting nil mask", t, func() {
		intial := Mask(errors.New("ABC"), errors.New("001"))
		err := Mask(intial, nil)
		_, ok := err.(*xerr)
		So(err.Error(), ShouldEqual, "ABC")
		So(ok, ShouldEqual, true)
	})
}

// TestGetDataAndSetData -
func TestGetDataAndSetData(t *testing.T) {
	Convey("GetData() and SetData()", t, func() {
		var val interface{}
		var ok bool

		val, ok = GetData(nil, "SOME_DATA")
		So(val, ShouldEqual, nil)
		So(ok, ShouldEqual, false)

		val, ok = GetData(errors.New("ABC"), "SOME_DATA")
		So(val, ShouldEqual, nil)
		So(ok, ShouldEqual, false)

		xerr := New("ERROR")

		val, ok = GetData(xerr, "SOME_DATA")
		So(val, ShouldEqual, nil)
		So(ok, ShouldEqual, false)

		SetData(xerr, "SOME_DATA", 100)

		val, ok = GetData(xerr, "SOME_DATA")
		So(val, ShouldEqual, 100)
		So(ok, ShouldEqual, true)
	})
}

// TestStack -
func TestStack(t *testing.T) {
	Convey("testing stack of nil error", t, func() {
		So(Stack(nil), ShouldEqual, nil)
	})

	Convey("testing stack of basic error", t, func() {
		So(Stack(errors.New("ABC")), ShouldEqual, nil)
	})

	Convey("testing stack of xerr", t, func() {
		err := New("ABC")
		stack := Stack(err)

		stack = transformStack(stack)

		So(stack, ShouldResemble, []StackLocation{
			StackLocation{
				Function: "xerrs.TestStack.func3",
				File:     "xerrs_test.go",
				Line:     414,
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
				Function: "xerrs.TestStack",
				File:     "xerrs_test.go",
				Line:     413,
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
		})
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
			Description:   "nil error",
			InputError:    nil,
			InputMask:     errors.New("MASK"),
			InputMaxStack: 100,
			Output:        ``,
		},
		TestCase{
			Description:   "basic",
			InputError:    errors.New("ERROR"),
			InputMask:     errors.New("MASK"),
			InputMaxStack: 100,
			Output: `
[ERROR] ERROR
[MASK ERROR] MASK
[STACK]:
xerrs.TestDetails.func1 [xerrs_test.go:619]
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
xerrs.TestDetails [xerrs_test.go:615]
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
xerrs.TestDetails.func1 [xerrs_test.go:619]
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
xerrs.TestDetails [xerrs_test.go:615]
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
xerrs.TestDetails.func1 [xerrs_test.go:621]
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
xerrs.TestDetails [xerrs_test.go:615]
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
xerrs.TestDetails.func1 [xerrs_test.go:619]
convey.parseAction.func1 [discovery.go:80]
convey.(*context).conveyInner [context.go:261]
convey.rootConvey.func1 [context.go:110]`,
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.Description, t, func() {
			var err error

			if testCase.InputMask != nil {
				err = Mask(testCase.InputError, testCase.InputMask)
			} else {
				err = Extend(testCase.InputError)
			}

			if x, ok := err.(*xerr); ok {
				x.stack = transformStack(x.stack)
			}

			So(Details(err, testCase.InputMaxStack), ShouldEqual, testCase.Output)
		})
	}

	Convey("Running stack on the basic Go error", t, func() {
		So(Details(errors.New("ABC"), 5), ShouldEqual, "ABC")
	})
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
