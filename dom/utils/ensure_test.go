package utils_test

import (
	"errors"
	"testing"
	. "thop/dom/utils"
)

func Test_Ensure(t *testing.T) {
	t.Run("passes if condition is true", func(t *testing.T) {
		// expect
		Ensure(true, "foo")
	})

	t.Run("panics if condition is false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// expect
		Ensure(false, "foo")
	})
}

func Test_EnsureWithErr(t *testing.T) {
	t.Run("passes if condition is true", func(t *testing.T) {
		// expect
		EnsureWithErr(true, errors.New("foo"))
	})

	t.Run("panics if condition is false", func(t *testing.T) {
		// given
		err := errors.New("foo")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			} else if r != err {
				t.Errorf("The code panicked with %v, expected %v", r, err)
			}
		}()

		// expect
		EnsureWithErr(false, err)
	})
}
