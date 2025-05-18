package repository

import (
	"fmt"
	"phopper/domain"
	"slices"
)

type LocalProjectRepository struct {}

func (r LocalProjectRepository) GetAllProjects() []domain.TmuxProject {
	return []domain.TmuxProject{
		{
			UUID: "1",
			Name: "Project 1",
			Path: "path 1",
		},
		{
			UUID: "2",
			Name: "Project 2",
			Path: "path 2",
		},
		{
			UUID: "3",
			Name: "Project 3",
			Path: "path 3",
		},
	}
}

func (r LocalProjectRepository) SaveProject(project domain.TmuxProject) domain.TmuxProject {
	project.UUID = fmt.Sprintf("%d", len(r.GetAllProjects()) + 1)
	return project
}

func (r LocalProjectRepository) DeleteProject(uuid string) {
	projects := r.GetAllProjects()
	for i, project := range projects {
		if project.UUID == uuid {
			fmt.Println("Deleting:", project)
			projects = slices.Delete(projects, i, i+1)
			break
		}
	}
}

