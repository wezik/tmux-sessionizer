package pane

import (
	"thop/dom/model/command"
	"thop/dom/problem"
)

type Name string
type Root string

const (
	ErrEmptyName = problem.Key("PANE_EMPTY_NAME")
)

// TODO: this structure is prepared but not used for now, panes need to be implemented properly
type Pane struct {
	Name     Name            `yaml:"name"`
	Root     Root            `yaml:"root,omitempty"`
	Commands command.Command `yaml:"run,omitempty"`
}

func New(name Name) *Pane {
	return &Pane{Name: name}
}

func (p *Pane) Validate() error {
	if p.Name == "" {
		return ErrEmptyName.WithMessage("name cannot be empty")
	}

	return nil
}
