package window

import (
	"thop/dom/model/command"
	"thop/dom/model/pane"
	"thop/dom/problem"
)

type Name string
type Root string
type Window struct {
	Name     Name              `yaml:"name"`
	Root     Root              `yaml:"root,omitempty"`
	Commands []command.Command `yaml:"run,omitempty"`
	Panes    []*pane.Pane      `yaml:"panes,omitempty"`
}

const (
	ErrEmptyName = problem.Key("WINDOW_EMPTY_NAME")
)

func New(name Name) *Window {
	return &Window{Name: name}
}

func (w *Window) Validate() error {
	if w.Name == "" {
		return ErrEmptyName.WithMessage("name cannot be empty")
	}

	// TODO: once panes are implemented, validate they're not empty here

	return nil
}
