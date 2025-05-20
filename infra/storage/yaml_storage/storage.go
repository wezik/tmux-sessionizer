package yaml_storage

import "phopper/domain/storage/repository"

type YamlStorage struct{}

func (y YamlStorage) GetProjectRepository() repository.ProjectRepository {
	return YamlProjectRepository{}
}

func (y YamlStorage) RunMigrations() {
	// TODO: implement
}
