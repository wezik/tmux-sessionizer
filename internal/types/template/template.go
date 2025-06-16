package template

import (
	"thop/internal/types/command"
	"thop/internal/types/window"
)

type Name string
type Root string
type ActiveWindow string

type Template struct {
	// Template name is used to specify the session name in multiplexer,
	// if not specified, the project name should be used
	Name         Name              `yaml:"name,omitempty"`
	Root         Root              `yaml:"root"`
	Commands     []command.Command `yaml:"run,omitempty"`
	Windows      []window.Window   `yaml:"windows"`
	ActiveWindow ActiveWindow      `yaml:"active_window,omitempty"`
}
