package domain

import (
	"fmt"
)

type TmuxProject struct {
	UUID string
	Name string
	Path string
}

func (project TmuxProject) String() string {
	return fmt.Sprintf("UUID: %s (Name: %s, Path: %s)", project.UUID, project.Name, project.Path)
}

type ProjectRepository interface {
	GetAllProjects() []TmuxProject
	SaveProject(TmuxProject) TmuxProject
}

type Config interface {
	UseFzf() bool
}
