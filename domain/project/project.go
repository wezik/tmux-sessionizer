package project

import "phopper/domain/project/template"

type Project struct {
	UUID    string            `yaml:"-"` // skip this field when marshalling
	Session template.Template `yaml:"session"`
}

func New(path string) Project {
	return Project{Session: template.New(path)}.WithDefaults()
}

// this function essentially sets defaults which allows to handle missing fields in case of template changes
func (p Project) WithDefaults() Project {
	p.Session = p.Session.WithDefaults()
	return p
}
