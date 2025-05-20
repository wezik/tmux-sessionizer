package storage

import "phopper/domain/storage/repository"

type Storage interface {
	GetProjectRepository() repository.ProjectRepository
}
