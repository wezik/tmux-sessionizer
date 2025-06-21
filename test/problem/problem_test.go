package problem_test

import (
	"errors"
	"testing"
	"thop/internal/problem"

	"github.com/stretchr/testify/assert"
)

func Test_WithMsg(t *testing.T) {
	t.Run("creates error with message", func(t *testing.T) {
		// given
		const key problem.Key = "test"
		const msg string = "test message"

		// when
		var err error = key.WithMsg(msg)

		// then
		assert.Equal(t, msg, err.Error())
		assert.Equal(t, key, err.(problem.Problem).Key)
	})
}

func Test_Equal(t *testing.T) {
	t.Run("returns true if keys are equal", func(t *testing.T) {
		// given
		const key problem.Key = "test"

		// when
		var err error = key.WithMsg("test message")

		// then
		assert.True(t, key.Equal(err))
	})

	t.Run("returns false if error is not a problem", func(t *testing.T) {
		// given
		const key problem.Key = "test"

		// when
		var err error = errors.New("some error")

		// then
		assert.False(t, key.Equal(err))
	})
}
