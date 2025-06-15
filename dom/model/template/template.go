package template

import (
	"thop/dom/model/command"
	"thop/dom/model/window"
	"thop/dom/problem"
)

type Name string
type Root string
type Template struct {
	// If desired, template name can be used to differentiate session name from the project name
	Name         Name              `yaml:"name,omitempty"`
	Root         Root              `yaml:"root"`
	Commands     []command.Command `yaml:"run,omitempty"`
	Windows      []*window.Window  `yaml:"windows"`
	ActiveWindow window.Name       `yaml:"active_window,omitempty"`
}

const (
	ErrNoWindows = problem.Key("TEMPLATE_NO_WINDOWS")
	ErrEmptyRoot = problem.Key("TEMPLATE_EMPTY_ROOT")
)

func New(root Root, windows []*window.Window) *Template {
	return &Template{
		Root:    root,
		Windows: windows,
	}
}

func (t *Template) Validate() error {
	if t.Root == "" {
		return ErrEmptyRoot.WithMessage("root cannot be empty")
	}

	if len(t.Windows) == 0 {
		return ErrNoWindows.WithMessage("at least one window is required")
	}

	return nil
}
