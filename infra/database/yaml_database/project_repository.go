package yaml_database

import (
	"os"
	"path/filepath"
	"phopper/domain/errors"
	"phopper/domain/project"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type YamlProjectRepository struct {}

func (y YamlProjectRepository) GetProjects() []project.Project {
	dir := getConfigPath()

	files, err := os.ReadDir(dir)
	errors.EnsureNotNil(err, "Could not read directory")

	projects := make([]project.Project, 0)

	for _, file := range files {
		if file.IsDir() {
			templateFile := filepath.Join(dir, file.Name(), "template.yaml")

			f, err := os.ReadFile(templateFile)
			errors.EnsureNotNil(err, "Could not read template file")

			var p project.Project
			yaml.Unmarshal(f, &p)

			p.UUID = file.Name()
			projects = append(projects, p)
		}
	}

	return projects
}

func (y YamlProjectRepository) SaveProject(project project.Project) project.Project {
	if project.UUID == "" {
		project.UUID = uuid.New().String()
	}

	dir := getConfigPath()
	templateFile := filepath.Join(dir, project.UUID, "template.yaml")

	err := os.MkdirAll(filepath.Dir(templateFile), 0755)
	errors.EnsureNotNil(err, "Could not create directory")

	f, err := os.Create(templateFile)
	errors.EnsureNotNil(err, "Could not create template file")
	defer f.Close()

	marshalled, err := yaml.Marshal(project)
	errors.EnsureNotNil(err, "Could not marshal project")

	_, err = f.Write(marshalled)
	errors.EnsureNotNil(err, "Could not write to template file")
	return project
}

func (y YamlProjectRepository) DeleteProject(uuid string) {
	dir := getConfigPath()
	err := os.RemoveAll(filepath.Join(dir, uuid))
	errors.EnsureNotNil(err, "Could not delete project")
}

func getConfigPath() string {
	cfg, err := os.UserConfigDir()
	errors.EnsureNotNil(err, "Could not get user config dir")
	return filepath.Join(cfg, ".phop")
}
