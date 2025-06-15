package service

import (
	"thop/dom/executor"
	"thop/dom/model/project"
	"thop/dom/model/template"
	"thop/dom/model/window"
	"thop/dom/multiplexer"
	"thop/dom/problem"
	"thop/dom/selector"
	"thop/dom/storage"
)

type Service interface {
	CreateProject(template.Root, project.Name) error
	DeleteProject(project.Name) error
	EditProject(project.Name) error
	OpenProject(project.Name) error
}

type ServiceImpl struct {
	selector       selector.Selector
	multiplexer    multiplexer.Multiplexer
	storage        storage.Storage
	executor executor.CommandExecutor
}

const (
	promptDeleteProject = "Select project to delete > "
	promptEditProject   = "Select project to edit > "
	promptOpenProject   = "Select project to open > "
)

const (
	ErrProjectNotFound problem.Key = "THOP_PROJECT_NOT_FOUND"
)

func New(
	selector selector.Selector,
	multiplexer multiplexer.Multiplexer,
	storage storage.Storage,
	executor executor.CommandExecutor,
) *ServiceImpl {
	return &ServiceImpl{
		selector: selector,
		multiplexer: multiplexer,
		storage: storage,
		executor: executor,
	}
}

func defaultProject(cwd template.Root, name project.Name) *project.Project {
	windows := []*window.Window{window.New("shell")}
	template := template.New(cwd, windows)
	return project.New(name, template)
}

func (s *ServiceImpl) CreateProject(cwd template.Root, name project.Name) error {
	p := defaultProject(cwd, name)

	// TODO: think on externalizing this validation block to a separate package
	// ****
	err := p.Validate()
	if err != nil {
		return err
	}

	p.Template.Validate()

	for _, w := range p.Template.Windows {
		w.Validate()
	}
	// ****

	return s.storage.Save(p)
}

func (s *ServiceImpl) OpenProject(name project.Name) error {
	project, err := s.resolveProjectName(name, promptOpenProject)
	if err != nil {
		return err
	}

	return s.multiplexer.AttachProject(project)
}

func (s *ServiceImpl) DeleteProject(name project.Name) error {
	project, err := s.resolveProjectName(name, promptDeleteProject)
	if err != nil {
		return err
	}

	return s.storage.Delete(project.UUID)
}

func (s *ServiceImpl) EditProject(name project.Name) error {
	project, err := s.resolveProjectName(name, promptEditProject)
	if err != nil {
		return err
	}

	templatePath, err := s.storage.PrepareTemplateFile(project)
	if err != nil {
		return err
	}

	cmd := executor.Command("nvim", templatePath)

	exitCode, err := s.executor.ExecuteInteractive(cmd)
	if err != nil {
		return err
	}

	if exitCode != executor.ExitCodeSuccess {
		return executor.ErrFailedExecution.WithMessage("program exited with non-zero exit code")
	}

	return nil
}

func (s *ServiceImpl) resolveProjectName(name project.Name, prompt string) (*project.Project, error) {
	if name != "" {
		return s.storage.Find(name)
	}

	projects, err := s.storage.List()
	if err != nil {
		return nil, err
	}

	entries := make([]string, len(projects))
	projectMap := make(map[string]*project.Project)

	for i, item := range projects {
		entries[i] = string(item.Name)
		projectMap[string(item.Name)] = item
	}

	selected, err := s.selector.SelectFrom(entries, prompt)
	if err != nil {
		return nil, err
	}

	if selected, ok := projectMap[selected]; ok {
		return selected, nil
	}

	return nil, ErrProjectNotFound.WithMessage(selected)
}
