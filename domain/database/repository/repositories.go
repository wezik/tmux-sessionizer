package repository

import "phopper/domain/project"

type ProjectRepository interface {
	GetProjects() []project.Project
	SaveProject(project.Project) project.Project
	DeleteProject(string)

	// decided to put this part of domain in the repository
	// this is really dependant on wether we use fs or db storage for the project templates
	// TODO: adjust implementation to make more sense domain wise
	PrepareTemplateFilePath(project.Project) string
}
