package utils

import "testing"

func Assert(t *testing.T, condition bool, message string, args ...any) {
	if !condition {
		t.Errorf(message, args...)
	}
}
