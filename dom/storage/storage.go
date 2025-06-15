package storage

import "thop/dom/model/project"

type Storage interface {
	List() ([]*project.Project, error)
	Find(name project.Name) (*project.Project, error)
	Save(t *project.Project) error
	Delete(uuid project.UUID) error
	PrepareTemplateFile(t *project.Project) (string, error)
}
