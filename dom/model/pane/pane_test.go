package pane_test

import (
	"testing"
	"thop/dom/model/pane"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	t.Run("passes for valid pane", func(t *testing.T) {
		// given
		p := pane.New("foobar")

		// expect
		assert.Nil(t, p.Validate())
	})

	t.Run("fails for empty name", func(t *testing.T) {
		// given
		p := pane.New("")

		// expect
		assert.Equal(t, pane.ErrEmptyName.WithMessage("name cannot be empty"), p.Validate())
	})
}
