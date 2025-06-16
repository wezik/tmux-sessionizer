package project

import (
	"thop/internal/types"
	"thop/internal/types/template"
)

type UUID string
type Name string

type Project struct {
	UUID     UUID              `yaml:"-"`
	Name     Name              `yaml:"name"`
	Version  types.Version     `yaml:"version"`
	Template template.Template `yaml:"template"`
}
