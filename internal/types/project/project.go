package project

import (
	"thop/internal/types"
	"thop/internal/types/template"
)

type UUID string
type Name string

type ProjectType int

const (
	TypeTemplate ProjectType = iota // will default to TypeTemplate if not set explicitly
	TypeTmuxSession
)

type Project struct {
	UUID     UUID              `yaml:"-"`
	Name     Name              `yaml:"name"`
	Version  types.Version     `yaml:"version"`
	Template template.Template `yaml:"template"`
	Type     ProjectType       `yaml:"-"`
}
