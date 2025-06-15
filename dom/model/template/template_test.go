package template_test

import (
	"testing"
	"thop/dom/model/template"
	"thop/dom/model/window"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	t.Run("passes for valid template", func(t *testing.T) {
		// given
		temp := template.New("some/root", []*window.Window{{}})

		// expect
		assert.Nil(t, temp.Validate())
	})

	t.Run("fails for no windows", func(t *testing.T) {
		// given
		temp := template.New("some/root", []*window.Window{})

		// expect
		assert.Equal(t, template.ErrNoWindows.WithMessage("at least one window is required"), temp.Validate())
	})

	t.Run("fails with empty root", func(t *testing.T) {
		// given
		temp := template.New("", []*window.Window{{}})

		// expect
		assert.Equal(t, template.ErrEmptyRoot.WithMessage("root cannot be empty"), temp.Validate())
	})
}
