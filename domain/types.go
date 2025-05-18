package domain

import (
	"fmt"
	"strings"
)

type TmuxProject struct {
	UUID string
	Name string
	Path string
}

func (project TmuxProject) String() string {
	return fmt.Sprintf("UUID:%s:Name:%s:Path:%s", project.UUID, project.Name, project.Path)
}

func TmuxProjectFromString(str string) TmuxProject {
	parts := strings.Split(str, ":")
	return TmuxProject{
		UUID: parts[1],
		Name: parts[3],
		Path: parts[5],
	}
}

type ProjectRepository interface {
	GetAllProjects() []TmuxProject
	SaveProject(TmuxProject) TmuxProject
	DeleteProject(string)
}

type Selector interface {
	ListAndSelect(entries []string, prompt string) string
}
