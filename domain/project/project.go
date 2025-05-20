package project

import "phopper/domain/project/session_template"

type Project struct {
	UUID string                              `yaml:"-"` // skip this field when marshalling
	Session session_template.SessionTemplate `yaml:"session"`
}

func New(path string) Project {
	return Project{ Session: session_template.New(path) }
}
