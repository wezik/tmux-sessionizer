package template

import (
	"os"
	"phopper/domain/errors"
)

const VERSION = 1

type Template struct {
	Version int    `yaml:"version"`
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
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

	return s
}
