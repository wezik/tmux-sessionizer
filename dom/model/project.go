package model

import (
	"errors"
	"fmt"
	. "thop/dom/utils"

	"github.com/google/uuid"
)

type Version int

const (
	V1 Version = 1
)

type Project struct {
	ID       string    `yaml:"-"`
	Name     string    `yaml:"name"`
	Version  Version   `yaml:"version"`
	Template *Template `yaml:"template"`
}

type Template struct {
	// Template name is used to specify the session name in multiplexer,
	// if not specified, the project name will be used
	Name         string   `yaml:"name,omitempty"`
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

func NewProject(name string, template *Template) (*Project, error) {
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

func NewTemplate(root string, windows []Window) (*Template, error) {
	EnsureWithErr(len(windows) > 0, ErrNeedsWindow)

	return &Template{
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

func (w Window) String() string {
	return fmt.Sprintf("{ Name: %s }", w.Name)
}

func (t Template) String() string {
	return fmt.Sprintf("{ Name: %s, Root: %s, Commands: %s, ActiveWindow: %s, Windows: %s }", t.Name, t.Root, t.Commands, t.ActiveWindow, t.Windows)
}

func (p Project) String() string {
	return fmt.Sprintf("{ ID: %s, Name: %s, Version: %d, Template: %s }", p.ID, p.Name, p.Version, p.Template)
}

var (
	ErrInvalidName       = errors.New("name cannot be empty")
	ErrNotFound          = errors.New("not found")
	ErrNeedsWindow       = errors.New("at least one window is required")
	ErrSelectorCancelled = errors.New("selector cancelled")
)
