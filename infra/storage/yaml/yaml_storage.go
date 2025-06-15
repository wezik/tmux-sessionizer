package yaml

import (
	"fmt"
	"path/filepath"
	"thop/dom/cfg"
	"thop/dom/fsystem"
	"thop/dom/model/project"
	"thop/dom/storage"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type YamlStorage struct {
	Config  *cfg.Config
	FSystem fsystem.FileSystem
}

const (
	templateFileName = "template.yaml"
	templatesDirName = "templates"
)

func New(config *cfg.Config, fileSystem fsystem.FileSystem) *YamlStorage {
	return &YamlStorage{Config: config, FSystem: fileSystem}
}

func (s *YamlStorage) List() ([]*project.Project, error) {
	cfgPath := s.Config.ConfigPath.String()
	templatesDir := filepath.Join(cfgPath, templatesDirName)

	// from what I understand, running os.Stat to check if a dir exists is not really providing
	// any benefits, and can also introduce weird edge cases, so instead just run mkdir everytime
	if err := s.FSystem.MkdirAll(templatesDir); err != nil {
		return nil, storage.ErrFailedToCreateConfigDir.WithMessage(err.Error())
	}

	dirs, err := s.FSystem.ReadDir(templatesDir)
	if err != nil {
		return nil, storage.ErrFailedToReadTemplateDir.WithMessage(err.Error())
	}

	var projects []*project.Project

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirName := dir.Name()

		templateFile := filepath.Join(templatesDir, dirName, templateFileName)
		bytes, err := s.FSystem.ReadFile(templateFile)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var p project.Project
		if err = yaml.Unmarshal(bytes, &p); err != nil {
			fmt.Println(err)
			continue
		}

		p.UUID = project.UUID(dirName)
		projects = append(projects, &p)
	}

	return projects, nil
}

func (s *YamlStorage) Find(name project.Name) (*project.Project, error) {
	projects, err := s.List()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == name {
			return project, nil
		}
	}

	return nil, storage.ErrTemplateNotFound.WithMessage(string(name))
}

func (s *YamlStorage) Save(p *project.Project) error {
	if p.UUID == "" {
		p.UUID = project.UUID(uuid.New().String())
	}

	cfgPath := s.Config.ConfigPath.String()
	templateDir := filepath.Join(cfgPath, templatesDirName, string(p.UUID))

	if err := s.FSystem.MkdirAll(templateDir); err != nil {
		return storage.ErrFailedToCreateTemplateDir.WithMessage(err.Error())
	}

	templateFile := filepath.Join(templateDir, templateFileName)
	bytes, err := yaml.Marshal(p)
	if err != nil {
		return storage.ErrFailedToMarshalTemplateFile.WithMessage(err.Error())
	}

	err = s.FSystem.WriteFile(templateFile, bytes)
	if err != nil {
		return storage.ErrFailedToWriteTemplateFile.WithMessage(err.Error())
	}

	return nil
}

func (s *YamlStorage) Delete(uuid project.UUID) error {
	cfgPath := s.Config.ConfigPath.String()
	templateDir := filepath.Join(cfgPath, templatesDirName, string(uuid))
	err := s.FSystem.RemoveAll(templateDir)
	if err != nil {
		return storage.ErrFailedToDeleteTemplate.WithMessage(err.Error())
	}

	return nil
}

func (s *YamlStorage) PrepareTemplateFile(p *project.Project) (string, error) {
	cfgPath := s.Config.ConfigPath.String()
	return filepath.Join(cfgPath, templatesDirName, string(p.UUID), templateFileName), nil
}
