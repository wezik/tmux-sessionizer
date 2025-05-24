package service

import (
	"errors"
	"fmt"
	"os"
	. "phopper/src/domain/model"
	. "phopper/src/domain/utils"
)

type Service interface {
	CreateProject(cwd, name string)
	SelectAndOpenProject(name string)
	DeleteProject(name string)
	EditProject(name string)
}

type ServiceImpl struct {
	sc Selector
	mu Multiplexer
	st Storage
}

func NewService(sc Selector, mu Multiplexer, st Storage) *ServiceImpl {
	return &ServiceImpl{sc: sc, mu: mu, st: st}
}

type Selector interface {
	SelectFrom(items []string) (string, error)
}

type Multiplexer interface {
	AttachProject(p *Project) error
}

type Storage interface {
	List() ([]*Project, error)
	Find(name string) (*Project, error)
	Save(t *Project) error
	Delete(uuid string) error
}

func (s *ServiceImpl) CreateProject(cwd, name string) {
	Ensure(name != "", "name cannot be empty")
	Ensure(cwd != "", "cwd cannot be empty")

	window, err := NewWindow("shell")
	EnsureWithErr(err == nil, err)

	template, err := NewTemplate(name, cwd, []Window{*window})
	EnsureWithErr(err == nil, err)

	project, err := NewProject(name, *template)
	EnsureWithErr(err == nil, err)

	err = s.st.Save(project)
	EnsureWithErr(err == nil, err)
}

func (s *ServiceImpl) SelectAndOpenProject(name string) {
	project := func() *Project {
		if name == "" {
			projects, err := s.st.List()
			EnsureWithErr(err == nil, err)

			selected, err := s.selectProject(projects)
			if err == ErrSelectorCancelled {
				fmt.Println("Cancelled")
				os.Exit(0)
			} else if err != nil {
				panic(err)
			}

			return selected
		} else {
			project, err := s.st.Find(name)
			EnsureWithErr(err == nil, err)

			return project
		}
	}()

	err := s.mu.AttachProject(project)
	EnsureWithErr(err == nil, err)
}

func (s *ServiceImpl) DeleteProject(name string) {
	panic("unimplemented")
}

func (s *ServiceImpl) EditProject(name string) {
	panic("unimplemented")
}

func (s *ServiceImpl) selectProject(items []*Project) (*Project, error) {
	itemsStringified := make([]string, len(items))
	itemsMap := make(map[string]*Project)
	for i, item := range items {
		itemsStringified[i] = item.Name
		itemsMap[item.Name] = item
	}

	selectedString, err := s.sc.SelectFrom(itemsStringified)
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

var (
	ErrSelectorCancelled = errors.New("Selector cancelled")
)
