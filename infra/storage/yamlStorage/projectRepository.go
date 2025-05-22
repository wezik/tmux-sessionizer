package yamlStorage

import (
	"os"
	"path/filepath"
	"phopper/domain/project"
	"phopper/domain/project/template"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type YamlProjectRepository struct{}

func NewYamlProjectRepository() *YamlProjectRepository {
	return &YamlProjectRepository{}
}

const templateFileName = "template.yaml"

func (y *YamlProjectRepository) GetProjects() ([]project.Project, error) {
	dir := getConfigPath()

	files, err := func() ([]os.DirEntry, error) {
		files, err := os.ReadDir(dir)

		// create and retry on fail
		if err != nil {
			path := getConfigPath()
			err := os.MkdirAll(path, 0755)
			// errors.EnsureNotNil(err, "Could not create config dir")
			if err != nil {
				return nil, err
			}

			files, err = os.ReadDir(dir)
		}

		return files, err
	}()

	if err != nil {
		return nil, err
	}

	var projects []project.Project

	for _, file := range files {
		if file.IsDir() {
			templateFile := filepath.Join(dir, file.Name(), templateFileName)

			f, err := os.ReadFile(templateFile)
			if err != nil {
				return nil, err
			}

			var t template.Template
			err = yaml.Unmarshal(f, &t)
			if err != nil {
				return nil, err
			}

			projects = append(projects, project.Project{UUID: file.Name(), Template: &t})
		}
	}

	return projects, nil
}

func (y *YamlProjectRepository) SaveProject(project *project.Project) (*project.Project, error) {
	if project.UUID == "" {
		project.UUID = uuid.New().String()
	}

	path := getConfigPath()
	projectDir := filepath.Join(path, project.UUID)

	err := os.MkdirAll(projectDir, 0755)
	// errors.EnsureNotNil(err, "Could not create config dir")
	if err != nil {
		return nil, err
	}

	templateFile := filepath.Join(projectDir, templateFileName)

	f, err := os.Create(templateFile)
	// errors.EnsureNotNil(err, "Could not create template file")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	marshalled, err := yaml.MarshalWithOptions(project.Template, yaml.IndentSequence(true))
	// errors.EnsureNotNil(err, "Could not marshal project")
	if err != nil {
		return nil, err
	}

	_, err = f.Write(marshalled)
	// errors.EnsureNotNil(err, "Could not write to template file")
	return project, err
}

func (y *YamlProjectRepository) DeleteProject(uuid string) error {
	dir := getConfigPath()
	err := os.RemoveAll(filepath.Join(dir, uuid))
	// errors.EnsureNotNil(err, "Could not delete project")
	return err
}
