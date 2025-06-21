package service

import (
	"os/exec"
	"thop/internal/config"
	"thop/internal/executor"
	"thop/internal/multiplexer"
	"thop/internal/problem"
	"thop/internal/selector"
	"thop/internal/storage"
	"thop/internal/types"
	"thop/internal/types/project"
	"thop/internal/types/template"
	"thop/internal/types/window"
)

type Service interface {
	CreateProject(template.Root, project.Name) error
	OpenProject(project.Name) error
	DeleteProject(project.Name) error
	EditProject(project.Name) error
}

type AppService struct {
	Selector    selector.ProjectSelector
	Multiplexer multiplexer.Multiplexer
	Storage     storage.Storage
	Config      *config.Config
	E           executor.CommandExecutor
}

const (
	ErrEditorNotSet        problem.Key = "THOP_EDITOR_NOT_SET"
	ErrEmptyProjectName    problem.Key = "THOP_EMPTY_PROJECT_NAME"
	ErrEmptyRootPath       problem.Key = "THOP_EMPTY_ROOT_PATH"
	ErrSelectedNonExisting problem.Key = "THOP_SELECTED_NON_EXISTING"
)

const (
	TemplateVersion = types.V1
)

func (s *AppService) CreateProject(root template.Root, name project.Name) error {
	if name == "" {
		return ErrEmptyProjectName.WithMsg("project name cannot be empty")
	}

	if root == "" {
		return ErrEmptyRootPath.WithMsg("root path cannot be empty")
	}

	p := project.Project{
		Name:    name,
		Version: TemplateVersion,
		Template: template.Template{
			Root: root,
			Windows: []window.Window{
				{Name: "shell"},
			},
		},
	}

	return s.Storage.Save(&p)
}

func (s *AppService) OpenProject(name project.Name) error {
	p, err := s.findOrSelect(name, "Select project to open > ", true)
	if err != nil {
		return err
	}

	return s.Multiplexer.AttachProject(p)
}

func (s *AppService) DeleteProject(name project.Name) error {
	p, err := s.findOrSelect(name, "Select project to delete > ", false)
	if err != nil {
		return err
	}

	return s.Storage.Delete(p.UUID)
}

func (s *AppService) EditProject(name project.Name) error {
	p, err := s.findOrSelect(name, "Select project to edit > ", false)
	if err != nil {
		return err
	}

	templatePath, err := s.Storage.PrepareTemplateFile(p)
	if err != nil {
		return err
	}

	editor := s.Config.GetEditor()
	if editor == "" {
		return ErrEditorNotSet.WithMsg("$EDITOR environment variable is not set")
	}

	cmd := exec.Command(editor, templatePath)
	_, err = s.E.ExecuteInteractive(cmd)
	return err
}

func (s *AppService) findOrSelect(name project.Name, prompt string, withActiveSessions bool) (project.Project, error) {
	if name != "" {
		p, err := s.Storage.Find(name)
		if err != nil {
			// if we can't find project, if enabled we can try to fallback to active session
			if withActiveSessions {
				if !storage.ErrProjectNotFound.Equal(err) {
					return project.Project{}, err
				}

				sessions, err := s.Multiplexer.ListActiveSessions()
				if err != nil {
					return project.Project{}, err
				}

				for _, session := range sessions {
					if session.Name == name {
						return session, nil
					}
				}
			} else {
				return p, err
			}
		}
		return p, nil
	}

	projects, err := s.Storage.List()
	if err != nil {
		return project.Project{}, err
	}

	if withActiveSessions {
		sessions, err := s.Multiplexer.ListActiveSessions()
		if err != nil {
			return project.Project{}, err
		}
		projects = append(projects, sessions...)
	}

	selected, err := s.Selector.SelectFrom(projects, prompt)

	if err != nil {
		return project.Project{}, err
	}

	return *selected, nil
}
