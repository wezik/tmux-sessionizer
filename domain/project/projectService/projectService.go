package projectService

import (
	"errors"
	"fmt"
	"phopper/domain/multiplexer"
	"phopper/domain/project"
	"phopper/domain/selector"
	"phopper/domain/shell"
	"phopper/domain/storage"
	"phopper/infra/multiplexer/tmuxMultiplexer"
	"phopper/infra/selector/fzfSelector"
	"phopper/infra/storage/yamlStorage"
)

type ProjectService interface {
	Create(root string, name string) (*project.Project, error)
	List() ([]project.Project, error)
	Delete(name string) (*project.Project, error)
	Open(name string) error
	Edit(name string, editor string) (*project.Project, error)
}

type ProjectServiceImpl struct {
	Storage     storage.Storage
	Selector    selector.Selector
	ShellRunner shell.Runner
	Multiplexer multiplexer.Multiplexer
}

func NewProjectService() *ProjectServiceImpl {
	return &ProjectServiceImpl{
		Storage:     yamlStorage.NewYamlStorage(),
		Selector:    fzfSelector.NewFZFSelector(),
		ShellRunner: shell.NewDefaultRunner(),
		Multiplexer: tmuxMultiplexer.NewTmux(),
	}
}

var (
	promptDeleteProject = "Select project to delete > "
	promptEditProject   = "Select project to edit > "
	promptOpenProject   = "Select project to open > "
)

func (ps *ProjectServiceImpl) Create(root string, name string) (*project.Project, error) {
	p := project.NewProject(root, name)

	repo := ps.Storage.GetProjectRepository()

	p, err := repo.SaveProject(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (ps *ProjectServiceImpl) List() ([]project.Project, error) {
	return ps.Storage.GetProjectRepository().GetProjects()
}

func (ps *ProjectServiceImpl) Delete(name string) (*project.Project, error) {
	p, err := func() (*project.Project, error) {
		if name != "" {
			return ps.findNamed(name)
		}
		return ps.interactiveSelect(promptDeleteProject)
	}()

	if err != nil {
		return nil, err
	}

	if p == nil {
		fmt.Println("Delete canceled")
		return nil, nil
	}

	return p, ps.Storage.GetProjectRepository().DeleteProject(p.UUID)
}

func (ps *ProjectServiceImpl) Open(name string) error {
	p, err := func() (*project.Project, error) {
		if name != "" {
			return ps.findNamed(name)
		}
		return ps.interactiveSelect(promptOpenProject)
	}()

	if err != nil {
		return err
	}

	if p == nil {
		fmt.Println("Open canceled")
		return nil
	}

	fmt.Println("Opening project", p.Template.Name)
	ps.Multiplexer.AssembleAndAttach(p)
	return nil
}

func (ps *ProjectServiceImpl) Edit(name string, editor string) (*project.Project, error) {
	p, err := func() (*project.Project, error) {
		if name != "" {
			return ps.findNamed(name)
		}
		return ps.interactiveSelect(promptEditProject)
	}()

	if err != nil {
		return nil, err
	}

	if p == nil {
		fmt.Println("Edit canceled")
		return nil, nil
	}

	templatePath, err := ps.Storage.PrepareTemplateFile(p)
	if err != nil {
		return nil, err
	}

	_, err, _ = ps.ShellRunner.RunInteractive(editor, templatePath)
	return p, nil
}

func (ps *ProjectServiceImpl) findNamed(name string) (*project.Project, error) {
	repo := ps.Storage.GetProjectRepository()

	pjs, err := repo.GetProjects()
	if err != nil {
		return nil, err
	}

	for _, p := range pjs {
		if p.Template.Name == name {
			return &p, nil
		}
	}

	return nil, errors.New("project not found")
}

func (ps *ProjectServiceImpl) interactiveSelect(prompt string) (*project.Project, error) {
	repo := ps.Storage.GetProjectRepository()

	pjs, err := repo.GetProjects()
	if err != nil {
		return nil, err
	}

	selected, err := ps.Selector.SelectProject(pjs, prompt)
	if err != nil {
		return nil, err
	}

	return selected, nil
}
