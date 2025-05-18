package repository

import (
	"fmt"
	"phopper/domain"
	"slices"
)

type LocalProjectRepository struct {
	projects []domain.TmuxProject
}

func NewLocalProjectRepository() LocalProjectRepository {
	return LocalProjectRepository{
		projects: []domain.TmuxProject{
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
		},
	}
}

func (r LocalProjectRepository) nextUUID() string {
	return fmt.Sprintf("%d", len(r.GetAllProjects()) + 1)
}

func (r LocalProjectRepository) GetAllProjects() []domain.TmuxProject {
	if (r.projects == nil) {
		return []domain.TmuxProject{}
	}

	return r.projects
}

func (r LocalProjectRepository) SaveProject(project domain.TmuxProject) domain.TmuxProject {
	if (r.projects == nil) {
		r.projects = []domain.TmuxProject{}
	}

	project.UUID = r.nextUUID()
	r.projects = append(r.projects, project)
	return project
}

func (r LocalProjectRepository) DeleteProject(uuid string) {
	if (r.projects == nil) {
		return
	}

	projects := r.GetAllProjects()
	for i, project := range projects {
		if project.UUID == uuid {
			fmt.Println("Deleting:", project)
			projects = slices.Delete(projects, i, i+1)
			break
		}
	}
}

