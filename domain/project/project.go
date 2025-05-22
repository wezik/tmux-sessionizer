package project

import "phopper/domain/project/template"

type Project struct {
	UUID     string
	Template *template.Template
}

func NewProject(root string, name string) *Project {
	return &Project{Template: template.NewTemplate(root, name)}
}
