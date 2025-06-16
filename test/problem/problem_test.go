package problem_test

import (
	"testing"
	"thop/problem"

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

func Test_EqualKey(t *testing.T) {
	t.Run("returns true if keys are equal", func(t *testing.T) {
		// given
		const key problem.Key = "test"

		// when
		var prob = key.WithMsg("test message")
		var err error = key.WithMsg(nil)

		// then
		assert.True(t, prob.EqualKey(err))
	})
}
