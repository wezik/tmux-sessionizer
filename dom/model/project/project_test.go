package project_test

import (
	"testing"
	"thop/dom/model/project"
	"thop/dom/model/template"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	t.Run("passes for valid project", func(t *testing.T) {
		// given
		p := project.New("foobar", &template.Template{})

		// expect
		assert.Nil(t, p.Validate())
	})

	t.Run("fails for empty name", func(t *testing.T) {
		// given
		p := project.New("", &template.Template{})

		// expect
		assert.Equal(t, project.ErrEmptyName.WithMessage("name cannot be empty"), p.Validate())
	})

	t.Run("fails for missing template", func(t *testing.T) {
		// given
		p := project.New("foobar", nil)

		// expect
		assert.Equal(t, project.ErrMissingTemplate.WithMessage("template cannot be missing"), p.Validate())
	})
}

func Test_New(t *testing.T) {
	t.Run("generates UUID", func(t *testing.T) {
		// given
		p := project.New("foobar", &template.Template{})

		// expect
		assert.NotNil(t, p.UUID)
	})
}
