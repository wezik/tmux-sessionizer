package repository

import "phopper/domain/project"

type ProjectRepository interface {
	GetProjects() []project.Project
	SaveProject(project.Project) project.Project
	DeleteProject(string)
}

