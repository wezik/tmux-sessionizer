package service

import (
	"fmt"
	"os"
	"os/exec"
	. "thop/dom/model"
	. "thop/dom/utils"

	"github.com/dsnet/try"
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

	window := try.E1(NewWindow("shell"))
	template := try.E1(NewTemplate(cwd, []Window{*window}))
	project := try.E1(NewProject(name, template))

	try.E(s.st.Save(project))
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

	try.E(s.mu.AttachProject(project))
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

	try.E(s.st.Delete(project.ID))
}

func (s *ServiceImpl) EditProject(name string) {
	project, err := s.findOrSelect(name, "Select project to edit > ")

	if err == ErrSelectorCancelled {
		fmt.Println("Edit operation cancelled")
		return
	} else if err != nil {
		panic(err)
	}

	templatePath := try.E1(s.st.PrepareTemplateFile(project))

	try.E(s.el.Open(templatePath))
}

func (s *ServiceImpl) findOrSelect(name string, prompt string) (p *Project, err error) {
	defer try.Handle(&err)

	if name != "" {
		return s.st.Find(name)
	}

	projects := try.E1(s.st.List())
	selected := try.E1(s.selectProject(projects, prompt))

	return selected, nil
}

func (s *ServiceImpl) selectProject(items []*Project, prompt string) (p *Project, err error) {
	defer try.Handle(&err)

	itemsStringified := make([]string, len(items))
	itemsMap := make(map[string]*Project)

	for i, item := range items {
		itemsStringified[i] = item.Name
		itemsMap[item.Name] = item
	}

	selectedString := try.E1(s.sl.SelectFrom(itemsStringified, prompt))

	selected, ok := itemsMap[selectedString]
	Ensure(ok, "selected item that does not exist")

	return selected, nil
}
