package pane

import "thop/internal/types/command"

type Name string
type Root string

type Pane struct {
	Name     Name              `yaml:"name"`
	Root     Root              `yaml:"root,omitempty"`
	Commands []command.Command `yaml:"run,omitempty"`
}
