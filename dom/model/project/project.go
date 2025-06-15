package project

import (
	"thop/dom/model/template"
	"thop/dom/problem"

	"github.com/google/uuid"
)

type UUID string
type Name string
type Version int

const (
	V1 Version = 1
)

const (
	ErrEmptyName       = problem.Key("PROJECT_EMPTY_NAME")
	ErrMissingTemplate = problem.Key("PROJECT_MISSING_TEMPLATE")
)

type Project struct {
	UUID     UUID               `yaml:"-"` // don't write this to yaml
	Name     Name               `yaml:"name"`
	Version  Version            `yaml:"version"`
	Template *template.Template `yaml:"template"`
}

func New(name Name, template *template.Template) *Project {
	uuid := UUID(uuid.New().String())
	version := V1
	return &Project{
		UUID:     uuid,
		Name:     name,
		Version:  version,
		Template: template,
	}
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return ErrEmptyName.WithMessage("name cannot be empty")
	}

	if p.Template == nil {
		return ErrMissingTemplate.WithMessage("template cannot be missing")
	}

	return nil
}
