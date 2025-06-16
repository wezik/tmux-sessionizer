package storage

import (
	"thop/cfg"
	"thop/internal/fsystem"
	"thop/internal/types/project"

	// "fmt"
	// "path/filepath"
	//
	// "github.com/dsnet/try"
	// "github.com/goccy/go-yaml"
	// "github.com/google/uuid"
)

type Storage interface {
	List() ([]project.Project, error)
	Find(project.Name) (project.Project, error)
	Save(*project.Project) error
	Delete(uuid project.UUID) error
	PrepareTemplateFile(project.Project) (string, error)
}

type YamlStorage struct {
	config cfg.Config
	fs     fsystem.FileSystem
}

const (
	templateFileName = "template.yaml"
	templatesDirName = "templates"
)

// func (s *YamlStorage) List() ([]project.Project, error) {
// 	defer try.Handle(&err)
//
// 	cfgDir := s.config.GetConfigDir()
//
// 	templatesDir := filepath.Join(cfgDir, templatesDirName)
//
// 	// from what I understand, running os.Stat to check if a dir exists is not really providing
// 	// any benefits, and can also introduce weird edge cases, so instead just run mkdir everytime
// 	try.E(s.fs.MkdirAll(templatesDir))
//
// 	dirs := try.E1(s.fs.ReadDir(templatesDir))
//
// 	var projects []*Project
//
// 	for _, dir := range dirs {
// 		if !dir.IsDir() {
// 			continue
// 		}
//
// 		dirName := dir.Name()
//
// 		templateFile := filepath.Join(templatesDir, dirName, templateFileName)
// 		bytes, err := s.fs.ReadFile(templateFile)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
//
// 		var project Project
// 		if err = yaml.Unmarshal(bytes, &project); err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
//
// 		project.ID = dirName
// 		projects = append(projects, &project)
// 	}
//
// 	return projects, nil
// }
//
// func (s *YamlStorage) Find(name string) (_ *Project, err error) {
// 	defer try.Handle(&err)
//
// 	projects := try.E1(s.List())
//
// 	for _, project := range projects {
// 		if project.Name == name {
// 			return project, nil
// 		}
// 	}
//
// 	return nil, ErrNotFound
// }
//
// func (s *YamlStorage) Save(p *Project) (err error) {
// 	defer try.Handle(&err)
//
// 	if p.ID == "" {
// 		p.ID = uuid.New().String()
// 	}
//
// 	cfgDir := s.config.GetConfigDir()
// 	templateDir := filepath.Join(cfgDir, templatesDirName, p.ID)
//
// 	try.E(s.fs.MkdirAll(templateDir))
//
// 	templateFile := filepath.Join(templateDir, templateFileName)
// 	bytes := try.E1(yaml.Marshal(p))
//
// 	return s.fs.WriteFile(templateFile, bytes)
// }
//
// func (s *YamlStorage) Delete(uuid string) error {
// 	cfgDir := s.config.GetConfigDir()
// 	templateDir := filepath.Join(cfgDir, templatesDirName, uuid)
// 	return s.fs.RemoveAll(templateDir)
// }
//
// func (s *YamlStorage) PrepareTemplateFile(p *Project) (string, error) {
// 	cfgDir := s.config.GetConfigDir()
// 	return filepath.Join(cfgDir, templatesDirName, p.ID, templateFileName), nil
// }
