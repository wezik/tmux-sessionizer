package problem_test

import (
	"errors"
	"testing"
	"thop/dom/problem"

	"github.com/stretchr/testify/assert"
)

func Test_WithMessage(t *testing.T) {
	t.Run("assembles an error", func(t *testing.T) {
		// given
		key := problem.Key("foo")
		message := "bar"

		expected := errors.New("problem ocurred foo: bar")

		// expect
		assert.Equal(t, expected.Error(), key.WithMessage(message).Error())
	})
}
