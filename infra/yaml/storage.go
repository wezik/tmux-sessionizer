package yaml

import (
	"fmt"
	"path/filepath"
	. "phopper/cfg"
	. "phopper/dom/model"
	. "phopper/dom/service"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type YamlStorage struct {
	config Config
	fs     FileSystem
}

func NewYamlStorage(config Config, fileSystem FileSystem) *YamlStorage {
	return &YamlStorage{config: config, fs: fileSystem}
}

var (
	templateFileName = "template.yaml"
	templatesDirName = "templates"
)

func (s *YamlStorage) List() ([]*Project, error) {
	cfgDir := s.config.GetConfigDir()

	templatesDir := filepath.Join(cfgDir, templatesDirName)

	// from what I understand, running os.Stat to check if a dir exists is not really providing
	// any benefits, and can also introduce weird edge cases, so instead just run mkdir everytime
	if err := s.fs.MkdirAll(templatesDir); err != nil {
		return nil, err
	}

	dirs, err := s.fs.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}

	var projects []*Project

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		templateFile := filepath.Join(templatesDir, dir.Name(), templateFileName)
		bytes, err := s.fs.ReadFile(templateFile)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var project Project
		if err = yaml.Unmarshal(bytes, &project); err != nil {
			fmt.Println(err)
			continue
		}

		project.ID = dir.Name()
		projects = append(projects, &project)
	}

	return projects, nil
}

func (s *YamlStorage) Find(name string) (*Project, error) {
	projects, err := s.List()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == name {
			return project, nil
		}
	}

	return nil, ErrNotFound
}

func (s *YamlStorage) Save(p *Project) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	cfgDir := s.config.GetConfigDir()
	templateDir := filepath.Join(cfgDir, templatesDirName, p.ID)

	if err := s.fs.MkdirAll(templateDir); err != nil {
		return err
	}

	templateFile := filepath.Join(templateDir, templateFileName)
	bytes, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	return s.fs.WriteFile(templateFile, bytes)
}

func (s *YamlStorage) Delete(uuid string) error {
	cfgDir := s.config.GetConfigDir()
	templateDir := filepath.Join(cfgDir, templatesDirName, uuid)
	return s.fs.RemoveAll(templateDir)
}

func (s *YamlStorage) PrepareTemplateFile(p *Project) (string, error) {
	cfgDir := s.config.GetConfigDir()
	return filepath.Join(cfgDir, templatesDirName, p.ID, templateFileName), nil
}
