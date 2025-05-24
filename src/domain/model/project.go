package model

import (
	"errors"

	"github.com/google/uuid"
)

type Version int

const (
	V1 Version = 1
)

type Project struct {
	ID       string   `yaml:"-"`
	Name     string   `yaml:"name"`
	Version  Version  `yaml:"version"`
	Template Template `yaml:"template"`
}

type Template struct {
	Name         string   `yaml:"name"`
	Root         string   `yaml:"root"`
	Commands     []string `yaml:"run,omitempty"`
	ActiveWindow string   `yaml:"active_window,omitempty"`
	Windows      []Window `yaml:"windows"`
}

type Window struct {
	Name     string   `yaml:"name"`
	Root     string   `yaml:"root,omitempty"`
	Commands []string `yaml:"run,omitempty"`
	Panes    []Pane   `yaml:"panes,omitempty"`
}

type Pane struct {
	Name     string   `yaml:"name"`
	Root     string   `yaml:"root,omitempty"`
	Commands []string `yaml:"run,omitempty"`
}

func NewProject(name string, template Template) (*Project, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	id := uuid.New().String()

	return &Project{
		ID:       id,
		Name:     name,
		Version:  V1,
		Template: template,
	}, nil
}

func NewTemplate(name string, root string, windows []Window) (*Template, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	if len(windows) == 0 {
		return nil, ErrNeedsWindow
	}

	return &Template{
		Name:    name,
		Root:    root,
		Windows: windows,
	}, nil
}

func NewWindow(name string) (*Window, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	return &Window{
		Name: name,
	}, nil
}

var (
	ErrInvalidName = errors.New("name cannot be empty")
	ErrNotFound    = errors.New("not found")
	ErrNeedsWindow = errors.New("at least one window is required")
)
