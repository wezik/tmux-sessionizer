package storage

import (
	"thop/dom/model/project"
	"thop/dom/problem"
)

type Storage interface {
	List() ([]*project.Project, error)
	Find(project.Name) (*project.Project, error)
	Save(*project.Project) error
	Delete(project.UUID) error
	PrepareTemplateFile(*project.Project) (string, error)
}

const (
	ErrFailedToCreateConfigDir     problem.Key = "STORAGE_FAILED_TO_CREATE_CONFIG_DIR"
	ErrFailedToCreateTemplateDir   problem.Key = "STORAGE_FAILED_TO_CREATE_TEMPLATE_DIR"
	ErrFailedToDeleteTemplate      problem.Key = "STORAGE_FAILED_TO_DELETE_TEMPLATE"
	ErrFailedToReadTemplateFile    problem.Key = "STORAGE_FAILED_TO_READ_TEMPLATE_FILE"
	ErrFailedToReadTemplateDir     problem.Key = "STORAGE_FAILED_TO_READ_TEMPLATE_DIR"
	ErrFailedToWriteTemplateFile   problem.Key = "STORAGE_FAILED_TO_WRITE_TEMPLATE_FILE"
	ErrFailedToMarshalTemplateFile problem.Key = "STORAGE_FAILED_TO_MARSHAL_TEMPLATE_FILE"

	ErrTemplateNotFound problem.Key = "STORAGE_TEMPLATE_NOT_FOUND"
)
