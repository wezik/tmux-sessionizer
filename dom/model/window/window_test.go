package window_test

import (
	"testing"
	"thop/dom/model/window"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	t.Run("passes for valid window", func(t *testing.T) {
		// given
		w := window.New("foobar")

		// expect
		assert.Nil(t, w.Validate())
	})

	t.Run("fails for empty name", func(t *testing.T) {
		// given
		w := window.New("")

		// expect
		assert.Equal(t, window.ErrEmptyName.WithMessage("name cannot be empty"), w.Validate())
	})
}
