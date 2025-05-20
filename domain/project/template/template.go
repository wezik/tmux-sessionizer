package template

import (
	"os"
	"phopper/domain/errors"
)

const VERSION = 1

type Template struct {
	Version      int      `yaml:"version"`
	Name         string   `yaml:"name"`
	Path         string   `yaml:"path"`
	ActiveWindow string   `yaml:"active_window"`
	OnStartHook  []string `yaml:"on_start"`
	Windows      []Window `yaml:"windows"`
}

func New(path string) Template {
	return Template{Path: path}.WithDefaults()
}

func (s Template) WithDefaults() Template {
	if s.Path == "" {
		path, err := os.UserHomeDir()
		if err != nil {
			errors.EnsureNotNil(err, "Could not get user home dir")
		}
		s.Path = path
	}

	if s.Name == "" {
		s.Name = s.Path
	}

	if s.Version != VERSION {
		s.Version = VERSION
	}

	if len(s.OnStartHook) == 0 {
		s.OnStartHook = make([]string, 1)
	}

	if len(s.Windows) == 0 {
		s.Windows = make([]Window, 1)
	}

	for i := range s.Windows {
		s.Windows[i] = s.Windows[i].WithDefaults()
	}

	if s.ActiveWindow == "" {
		s.ActiveWindow = s.Windows[0].Name
	}

	return s
}

type Window struct {
	Name        string   `yaml:"name"`
	Path        string   `yaml:"path"`
	OnStartHook []string `yaml:"on_start"`
}

func (w Window) WithDefaults() Window {
	// name empty is fine
	// path empty is fine
	if len(w.OnStartHook) == 0 {
		w.OnStartHook = make([]string, 1)
	}
	return w
}
