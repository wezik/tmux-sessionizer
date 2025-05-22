package storage

import "phopper/domain/project"

type ProjectRepository interface {
	GetProjects() ([]project.Project, error)
	SaveProject(*project.Project) (*project.Project, error)
	DeleteProject(uuid string) error
}
