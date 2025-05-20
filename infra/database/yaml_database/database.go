package yaml_database

import "phopper/domain/database/repository"

type YamlDatabase struct {}

func (y YamlDatabase) GetProjectRepository() repository.ProjectRepository {
	return YamlProjectRepository{}
}

func (y YamlDatabase) RunMigrations() {
	// TODO: implement
}
