package service

import (
	"fmt"
	"os"
	"os/exec"
	. "thop/dom/model"
	. "thop/dom/utils"
)

type Service interface {
	CreateProject(cwd, name string)
	SelectAndOpenProject(name string)
	DeleteProject(name string)
	EditProject(name string)
}

type ServiceImpl struct {
	sl Selector
	mu Multiplexer
	st Storage
	el EditorLauncher
}

func NewService(sl Selector, mu Multiplexer, st Storage, el EditorLauncher) *ServiceImpl {
	return &ServiceImpl{sl: sl, mu: mu, st: st, el: el}
}

type Selector interface {
	SelectFrom(items []string, prompt string) (string, error)
}

type Multiplexer interface {
	AttachProject(p *Project) error
}

type CommandExecutor interface {
	Execute(cmd *exec.Cmd) (string, int, error)
	ExecuteInteractive(cmd *exec.Cmd) (int, error)
}

type Storage interface {
	List() ([]*Project, error)
	Find(name string) (*Project, error)
	Save(t *Project) error
	Delete(uuid string) error
	PrepareTemplateFile(t *Project) (string, error)
}

type FileSystem interface {
	MkdirAll(path string) error
	ReadDir(path string) ([]os.DirEntry, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	RemoveAll(path string) error
}

type EditorLauncher interface {
	Open(path string) error
}

func (s *ServiceImpl) CreateProject(cwd, name string) {
	Ensure(name != "", "name cannot be empty")
	Ensure(cwd != "", "cwd cannot be empty")

	window, err := NewWindow("shell")
	EnsureWithErr(err == nil, err)

	template, err := NewTemplate(cwd, []Window{*window})
	EnsureWithErr(err == nil, err)

	project, err := NewProject(name, template)
	EnsureWithErr(err == nil, err)

	err = s.st.Save(project)
	EnsureWithErr(err == nil, err)
}

func (s *ServiceImpl) SelectAndOpenProject(name string) {
	project, err := s.findOrSelect(name, "Select project to open > ")

	if err != nil {
		if err == ErrSelectorCancelled {
			fmt.Println("Select operation cancelled")
			return
		}
		panic(err)
	}

	err = s.mu.AttachProject(project)
	EnsureWithErr(err == nil, err)
}

func (s *ServiceImpl) DeleteProject(name string) {
	project, err := s.findOrSelect(name, "Select project to delete > ")

	if err != nil {
		if err == ErrSelectorCancelled {
			fmt.Println("Delete operation cancelled")
			return
		}
		panic(err)
	}

	err = s.st.Delete(project.ID)
	EnsureWithErr(err == nil, err)
}

func (s *ServiceImpl) EditProject(name string) {
	project, err := s.findOrSelect(name, "Select project to edit > ")
	if err == ErrSelectorCancelled {
		fmt.Println("Edit operation cancelled")
		return
	} else if err != nil {
		panic(err)
	}

	templatePath, err := s.st.PrepareTemplateFile(project)
	EnsureWithErr(err == nil, err)

	err = s.el.Open(templatePath)
	EnsureWithErr(err == nil, err)
}

func (s *ServiceImpl) findOrSelect(name string, prompt string) (*Project, error) {
	if name != "" {
		return s.st.Find(name)
	}

	projects, err := s.st.List()
	if err != nil {
		return nil, err
	}

	selected, err := s.selectProject(projects, prompt)
	if err == ErrSelectorCancelled {
		return nil, err
	}

	Ensure(err == nil, "Unknown error occured while selecting the project")

	return selected, nil
}

func (s *ServiceImpl) selectProject(items []*Project, prompt string) (*Project, error) {
	itemsStringified := make([]string, len(items))
	itemsMap := make(map[string]*Project)

	for i, item := range items {
		itemsStringified[i] = item.Name
		itemsMap[item.Name] = item
	}

	selectedString, err := s.sl.SelectFrom(itemsStringified, prompt)
	if err == ErrSelectorCancelled {
		return nil, err // in case of cancellation, propagate the error upwards
	}

	EnsureWithErr(err == nil, err)

	// this is a bit redundant, but the more fail-safes the better,
	// it would require faulty implementation of the selector
	EnsureWithErr(selectedString != "", ErrSelectorCancelled)

	selected, ok := itemsMap[selectedString]
	Ensure(ok, "selected item that does not exist")

	return selected, nil
}
