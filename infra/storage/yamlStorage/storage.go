package yamlStorage

import (
	"os"
	"path/filepath"
	"phopper/domain/errors"
	"phopper/domain/project"
	"phopper/domain/storage"
)

type YamlStorage struct {
	ProjectRepository storage.ProjectRepository
}

func NewYamlStorage() *YamlStorage {
	return &YamlStorage{
		ProjectRepository: NewYamlProjectRepository(),
	}
}

func (y *YamlStorage) GetProjectRepository() storage.ProjectRepository {
	return y.ProjectRepository
}

func (y *YamlStorage) PrepareTemplateFile(p *project.Project) (string, error) {
	return filepath.Join(getConfigPath(), p.UUID, templateFileName), nil
}

func getConfigPath() string {
	cfg, err := os.UserConfigDir()
	errors.EnsureNotNil(err, "Could not get user config dir")
	return filepath.Join(cfg, ".phop", "templates")
}
