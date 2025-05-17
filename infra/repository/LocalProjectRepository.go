package repository

import (
	"fmt"
	"phopper/domain"
)

type LocalProjectRepository struct {}

func (r LocalProjectRepository) GetAllProjects() []domain.TmuxProject {
	return []domain.TmuxProject{
		{
			Name: "Project 1",
			Path: "path 1",
		},
		{
			Name: "Project 2",
			Path: "path 2",
		},
		{
			Name: "Project 3",
			Path: "path 3",
		},
	}
}

func (r LocalProjectRepository) SaveProject(project domain.TmuxProject) domain.TmuxProject {
	project.UUID = fmt.Sprintf("%d", len(r.GetAllProjects()))
	return project
}

