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
	return fmt.Sprintf("UUID:%s:Name:%s:Path:%s", project.UUID, project.Name, project.Path)
}

type Database interface {
	RunMigrations()
}

type ProjectRepository interface {
	GetProjects() []TmuxProject
	SaveProject(project TmuxProject) TmuxProject
	DeleteProject(uuid string)
}

type Selector interface {
	ListAndSelect(entries []string, prompt string) string
}
