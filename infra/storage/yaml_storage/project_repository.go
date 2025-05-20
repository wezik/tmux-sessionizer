package yaml_storage

import (
	"fmt"
	"os"
	"path/filepath"
	"phopper/domain/errors"
	"phopper/domain/project"
	"phopper/domain/project/template"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type YamlProjectRepository struct{}

const templateFileName = "template.yaml"

func (y YamlProjectRepository) GetProjects() []project.Project {
	dir := getConfigPath()

	files, err := func() ([]os.DirEntry, error) {
		files, err := os.ReadDir(dir)

		// create and retry on fail
		if err != nil {
			createConfigDir()
			files, err = os.ReadDir(dir)
		}

		return files, err
	}()
	errors.EnsureNotNil(err, "Could not read directory")

	projects := make([]project.Project, 0)

	migrated := 0
	for _, file := range files {
		if file.IsDir() {
			templateFile := filepath.Join(dir, file.Name(), templateFileName)

			f, err := os.ReadFile(templateFile)
			errors.EnsureNotNil(err, "Could not read template file")

			var p project.Project
			yaml.Unmarshal(f, &p)

			p.UUID = file.Name()

			// on version mismatch, save with defaults before adding to projects
			if p.Session.Version != template.VERSION {
				p = y.SaveProject(p.WithDefaults())
				migrated++
			}
			projects = append(projects, p)
		}
	}

	if migrated > 0 {
		fmt.Println("Migrated", migrated, "templates to version", template.VERSION)
	}

	return projects
}

func (y YamlProjectRepository) SaveProject(project project.Project) project.Project {
	if project.UUID == "" {
		project.UUID = uuid.New().String()
	}

	dir := getConfigPath()
	templateFile := filepath.Join(dir, project.UUID, templateFileName)

	createConfigDir()

	f, err := os.Create(templateFile)
	errors.EnsureNotNil(err, "Could not create template file")
	defer f.Close()

	marshalled, err := yaml.Marshal(project.WithDefaults())
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

func (y YamlProjectRepository) PrepareTemplateFilePath(p project.Project) string {
	return filepath.Join(getConfigPath(), p.UUID, templateFileName)
}

func getConfigPath() string {
	cfg, err := os.UserConfigDir()
	errors.EnsureNotNil(err, "Could not get user config dir")
	return filepath.Join(cfg, ".phop", "templates")
}

func createConfigDir() {
	path := getConfigPath()
	err := os.MkdirAll(path, 0755)
	errors.EnsureNotNil(err, "Could not create config dir")
}
