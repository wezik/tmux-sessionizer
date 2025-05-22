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

func NewTemplate(root string, name string) *Template {
	if root == "" {
		if home, err := os.UserHomeDir(); err == nil {
			root = home
		} else {
			errors.EnsureNotNil(err, "Could not get user home dir")
		}
	}

	if name == "" {
		name = root
	}

	return &Template{
		Version: VERSION,
		Name:    name,
		Root:    root,
		Windows: []Window{NewWindow("default")},
	}
}

type Window struct {
	Name     string          `yaml:"name"`
	Root     string          `yaml:"root,omitempty"`
	Commands []string        `yaml:"run,omitempty"`
	Panes    map[string]Pane `yaml:"panes,omitempty"`
}

func NewWindow(name string) Window {
	return Window{Name: name}
}

type Pane struct {
	Name     string   `yaml:"name"`
	Root     string   `yaml:"root,omitempty"`
	Commands []string `yaml:"run,omitempty"`
}

func NewPane(name string) Pane {
	return Pane{Name: name}
}
