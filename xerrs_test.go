package xerrs

import (
	"errors"
	"strings"
	"testing"
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

func TestNew(t *testing.T) {
	in := New("ABC")

	if x, ok := in.(*xerr); ok {
		if x.cause.Error() != "ABC" {
			t.Errorf("wrong error cause: want=%v got=%v", "ABC", x.cause.Error())
		}
		if x.mask != nil {
			t.Errorf("expected nil mask, got=%v", x.mask)
		}
		if in.Error() != "ABC" {
			t.Errorf("wrong error message: want=%v got=%v", "ABC", in.Error())
		}
		if len(x.stack) != 3 {
			t.Errorf("wrong stack length: want=%v got=%v", 3, len(x.stack))
		}
	}
}

func TestErrorf(t *testing.T) {
	in := Errorf("some error %d %v", 1, "HELLO")
	if x, ok := in.(*xerr); ok {
		if x.cause.Error() != "some error 1 HELLO" {
			t.Errorf("wrong error cause: want=%v got=%v", "some error 1 HELLO", x.cause.Error())
		}
		if x.mask != nil {
			t.Errorf("expected nil mask: got=%v", x.mask)
		}
		if in.Error() != "some error 1 HELLO" {
			t.Errorf("wrong error message: want=%v got=%v", "some error 1 HELLO", in.Error())
		}
		if len(x.stack) != 3 {
			t.Errorf("wrong stack length: want=%v got=%v", 3, len(x.stack))
		}
	}
}

func TestExtend(t *testing.T) {
	in := Extend(errors.New("ABC"))

	if x, ok := in.(*xerr); ok {
		if x.cause.Error() != "ABC" {
			t.Errorf("wrong error cause: want=%v got=%v", "ABC", x.cause.Error())
		}
		if x.mask != nil {
			t.Errorf("expected nil mask, got=%v", x.mask)
		}
		if in.Error() != "ABC" {
			t.Errorf("wrong error message: want=%v got=%v", "ABC", in.Error())
		}
		if len(x.stack) != 3 {
			t.Errorf("wrong stack length: want=%v got=%v", 3, len(x.stack))
		}
	}
}

func TestMask(t *testing.T) {
	err := Mask(nil, errors.New("ABC"))
	if err != nil {
		t.Errorf("expected nil error: got=%v", err)
	}

	err = Mask(errors.New("ABC"), nil)
	_, ok := err.(*xerr)
	if err.Error() != "ABC" {
		t.Errorf("wrong error message: want=%v got=%v", "ABC", err.Error())
	}
	if ok != true {
		t.Errorf("expected err to be xerr")
	}

	err = Mask(errors.New("ABC"), errors.New("XYZ"))
	_, ok = err.(*xerr)
	if err.Error() != "XYZ" {
		t.Errorf("wrong error message: want=%v got=%v", "XYZ", err.Error())
	}
	if ok != true {
		t.Errorf("expected err to be xerr")
	}

	intial := Extend(errors.New("ABC"))
	err = Mask(intial, errors.New("XYZ"))
	_, ok = err.(*xerr)
	if err.Error() != "XYZ" {
		t.Errorf("wrong error message: want=%v got=%v", "XYZ", err.Error())
	}
	if ok != true {
		t.Errorf("expected err to be xerr")
	}

	intial = Mask(errors.New("ABC"), errors.New("001"))
	err = Mask(intial, errors.New("XYZ"))
	_, ok = err.(*xerr)
	if err.Error() != "XYZ" {
		t.Errorf("wrong error message: want=%v got=%v", "XYZ", err.Error())
	}
	if ok != true {
		t.Errorf("expected err to be xerr")
	}

	intial = Mask(errors.New("ABC"), errors.New("001"))
	err = Mask(intial, nil)
	_, ok = err.(*xerr)
	if err.Error() != "ABC" {
		t.Errorf("wrong error message: want=%v got=%v", "ABC", err.Error())
	}
	if ok != true {
		t.Errorf("expected err to be xerr")
	}
}

func TestData(t *testing.T) {
	t.Run("SetGet", func(t *testing.T) {
		err := New("test")
		SetData(err, "SOME_DATA", "test")

		v, ok := GetData(err, "SOME_DATA")
		if !ok {
			t.Error("expected data")
		}

		if _, ok := v.(string); !ok {
			t.Errorf("expected string, got %T", v)
		}
	})

	t.Run("nil error nil data", func(t *testing.T) {
		_, ok := GetData(nil, "test")
		if ok {
			t.Error("expected false")
		}
	})

	t.Run("nil data", func(t *testing.T) {
		err := New("test")
		if _, ok := GetData(err, "test"); ok {
			t.Error("expected false")
		}
	})

	t.Run("nested error data", func(t *testing.T) {
		a := New("a")
		SetData(a, "foo", "bar")

		b := Wrap(a, "b")

		foo, ok := GetData(b, "foo")
		if !ok {
			t.Error("failed to get data for \"foo\"")
		}

		msg, ok := foo.(string)
		if !ok {
			t.Error("data for \"foo\" should be a string")
		}

		if msg != "bar" {
			t.Errorf("wanted foo=%q, got foo=%q", "bar", msg)
		}
	})
}

func TestDetails(t *testing.T) {
	type TestCase struct {
		Description   string
		InputError    error
		InputMask     error
		InputMaxStack int
		OutputPrefix  string
	}

	testCases := []TestCase{
		TestCase{
			Description:   "nil error",
			InputError:    nil,
			InputMask:     errors.New("MASK"),
			InputMaxStack: 100,
			OutputPrefix:  ``,
		},
		TestCase{
			Description:   "basic",
			InputError:    errors.New("ERROR"),
			InputMask:     errors.New("MASK"),
			InputMaxStack: 100,
			OutputPrefix: `
[ERROR] ERROR
[MASK ERROR] MASK
[STACK]:`,
		},
		TestCase{
			Description:   "mask is the same as error",
			InputError:    errors.New("ERROR"),
			InputMask:     errors.New("ERROR"),
			InputMaxStack: 100,
			OutputPrefix: `
[ERROR] ERROR
[STACK]:`,
		},
		TestCase{
			Description:   "mask is nil",
			InputError:    errors.New("ERROR"),
			InputMask:     nil,
			InputMaxStack: 100,
			OutputPrefix: `
[ERROR] ERROR
[STACK]:`,
		},
		TestCase{
			Description:   "fewer stack lines",
			InputError:    errors.New("ERROR"),
			InputMask:     errors.New("MASK"),
			InputMaxStack: 4,
			OutputPrefix: `
[ERROR] ERROR
[MASK ERROR] MASK
[STACK]:`,
		},
	}

	for _, testCase := range testCases {
		var err error

		if testCase.InputMask != nil {
			err = Mask(testCase.InputError, testCase.InputMask)
		} else {
			err = Extend(testCase.InputError)
		}

		if x, ok := err.(*xerr); ok {
			x.stack = transformStack(x.stack)
		}

		if !strings.HasPrefix(Details(err, testCase.InputMaxStack), testCase.OutputPrefix) {
			t.Errorf("wrong output prefix: wanted prefix=%v got=%v", testCase.OutputPrefix, Details(err, testCase.InputMaxStack))
		}
	}

	if Details(errors.New("ABC"), 5) != "ABC" {
		t.Errorf("running stack on the basic go error failed: want=%v got=%v", "ABC", Details(errors.New("ABC"), 5))
	}
}

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
		t.Run(testCase.Description, func(t *testing.T) {
			if IsEqual(testCase.InputErr1, testCase.InputErr2) != testCase.Output {
				t.Errorf("wrong output: want=%v got=%v", testCase.Output, IsEqual(testCase.InputErr1, testCase.InputErr2))
			}
		})
	}
}

func TestWrap(t *testing.T) {
	for _, test := range []struct {
		in      error
		msg     string
		wantMsg string
		isNil   bool
	}{
		{
			in:      nil,
			msg:     "nothing",
			wantMsg: "",
			isNil:   true,
		},
		{
			in:      New("i/o error"),
			msg:     "read",
			wantMsg: "read: i/o error",
			isNil:   false,
		},
	} {
		err := Wrap(test.in, test.msg)
		if err == nil && !test.isNil {
			t.Error("expected non-nil error")
		} else if err != nil && test.isNil {
			t.Errorf("expected nil error, got=%v", err)
		}

		if err != nil {
			if got := err.Error(); test.wantMsg != got {
				t.Errorf("wrong error message: want=%v got=%v", test.wantMsg, got)
			}
		}
	}
}

func TestWrapf(t *testing.T) {
	for _, test := range []struct {
		in      error
		format  string
		args    []interface{}
		wantMsg string
		isNil   bool
	}{
		{
			in:      nil,
			format:  "nothing",
			args:    nil,
			wantMsg: "",
			isNil:   true,
		},
		{
			in:      New("i/o error"),
			format:  "read %q",
			args:    []interface{}{"config.yaml"},
			wantMsg: "read \"config.yaml\": i/o error",
			isNil:   false,
		},
	} {
		err := Wrapf(test.in, test.format, test.args...)
		if err == nil && !test.isNil {
			t.Error("expected non-nil error")
		} else if err != nil && test.isNil {
			t.Errorf("expected nil error, got=%v", err)
		}

		if err != nil {
			if got := err.Error(); test.wantMsg != got {
				t.Errorf("wrong error message: want=%v got=%v", test.wantMsg, got)
			}
		}
	}
}

func BenchmarkWrap(b *testing.B) {
	b.Run("chain", func(b *testing.B) {
		err := New("test")
		for i := 0; i < b.N; i++ {
			wrapped := &xerr{
				stack: getStack(stackFunctionOffset),
				cause: err,
				msg:   "prepend",
			}

			// inlined (*xerr).Error
			if wrapped.mask != nil {
				_ = wrapped.mask.Error()
				continue
			}

			if wrapped.msg != "" {
				_ = wrapped.msg + ": " + wrapped.cause.Error()
				continue
			}

			_ = wrapped.cause.Error()
		}
	})

	b.Run("prepend", func(b *testing.B) {
		err := New("test")
		for i := 0; i < b.N; i++ {
			wrapped := &xerr{
				stack: getStack(stackFunctionOffset),
				cause: err,
				msg:   "prepend" + ": " + err.Error(),
			}

			// inlined (*xerr).Error
			if wrapped.mask != nil {
				_ = wrapped.mask.Error()
				continue
			}

			if wrapped.msg != "" {
				_ = wrapped.msg + ": " + wrapped.cause.Error()
				continue
			}

			_ = wrapped.cause.Error()
		}
	})

	b.Run("prepend with typecheck", func(b *testing.B) {
		for name, err := range map[string]error{
			"error": errors.New("test"),
			"xerr":  New("test"),
		} {
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if xe, ok := err.(*xerr); ok {
						xe.msg = "prepend" + ": " + xe.msg

						// modified (*xerr).Error
						if xe.mask != nil {
							_ = xe.mask.Error()
							continue
						}

						if xe.msg != "" {
							continue
						}

						_ = xe.cause.Error()
					} else {
						wrapped := &xerr{
							stack: getStack(stackFunctionOffset),
							cause: err,
							msg:   "prepend" + ": " + err.Error(),
						}

						// modified (*xerr).Error
						if wrapped.mask != nil {
							_ = wrapped.mask.Error()
							continue
						}

						if wrapped.msg != "" {
							continue
						}

						_ = wrapped.cause.Error()
					}
				}
			})
		}
	})
}
