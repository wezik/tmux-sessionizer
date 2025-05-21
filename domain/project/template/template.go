package template

import (
	"os"
	"phopper/domain/errors"
)

const VERSION = 1

type Template struct {
	Version      int      `yaml:"version"`
	Name         string   `yaml:"name"`
	Root         string   `yaml:"root"`
	Commands     []string `yaml:"run,omitempty"`
	ActiveWindow string   `yaml:"active_window,omitempty"`
	Windows      []Window `yaml:"windows"`
}

func New(path string) Template {
	return Template{Root: path}.WithDefaults()
}

func (s Template) WithDefaults() Template {
	if s.Root == "" {
		path, err := os.UserHomeDir()
		if err != nil {
			errors.EnsureNotNil(err, "Could not get user home dir")
		}
		s.Root = path
	}

	if s.Name == "" {
		s.Name = s.Root
	}

	if s.Version != VERSION {
		s.Version = VERSION
	}

	if len(s.Windows) == 0 {
		s.Windows = append(s.Windows, Window{Name: "default"})
	}

	for i := range s.Windows {
		s.Windows[i] = s.Windows[i].WithDefaults()
	}

	return s
}

type Window struct {
	Name     string          `yaml:"name"`
	Root     string          `yaml:"root,omitempty"`
	Commands []string        `yaml:"run,omitempty"`
	Panes    map[string]Pane `yaml:"panes,omitempty"`
}

func (w Window) WithDefaults() Window {
	for name, pane := range w.Panes {
		w.Panes[name] = pane.WithDefaults()
	}
	return w
}

type Pane struct {
	Name     string   `yaml:"name"`
	Root     string   `yaml:"root,omitempty"`
	Commands []string `yaml:"run,omitempty"`
}

func (p Pane) WithDefaults() Pane {
	return p
}
