package storage

import (
	"phopper/domain/project"
)

type Storage interface {
	GetProjectRepository() ProjectRepository
	PrepareTemplateFile(p *project.Project) (string, error)
}
