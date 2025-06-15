package service_test

import (
	"fmt"
	"testing"
	. "thop/dom/model"
	. "thop/dom/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSelector struct {
	mock.Mock
}

func (s *MockSelector) SelectFrom(items []string, prompt string) (string, error) {
	args := s.Called(items, prompt)
	return args.String(0), args.Error(1)
}

type MockMultiplexer struct {
	mock.Mock
}

func (s *MockMultiplexer) AttachProject(p *Project) error {
	args := s.Called(p)
	return args.Error(0)
}

type MockStorage struct {
	mock.Mock
}

func (s *MockStorage) List() ([]*Project, error) {
	args := s.Called()
	return args.Get(0).([]*Project), args.Error(1)
}

func (s *MockStorage) Find(name string) (*Project, error) {
	args := s.Called(name)
	return args.Get(0).(*Project), args.Error(1)
}

func (s *MockStorage) Save(t *Project) error {
	args := s.Called(t)
	return args.Error(0)
}

func (s *MockStorage) Delete(uuid string) error {
	args := s.Called(uuid)
	return args.Error(0)
}

func (s *MockStorage) PrepareTemplateFile(t *Project) (string, error) {
	args := s.Called(t)
	return args.String(0), args.Error(1)
}

type MockEditorLauncher struct {
	mock.Mock
}

func (s *MockEditorLauncher) Open(path string) error {
	args := s.Called(path)
	return args.Error(0)
}

func Test_CreateProject(t *testing.T) {
	t.Run("creates project", func(t *testing.T) {
		// given
		stMock := new(MockStorage)
		svc := New(nil, nil, stMock, nil)

		cwd := "/home/test"
		name := "foobar"

		var projectParam *Project
		stMock.On("Save", mock.Anything).Run(func(args mock.Arguments) {
			projectParam = args.Get(0).(*Project)
		}).Return(nil)

		// when
		svc.CreateProject(cwd, name)

		// then
		assert.Equal(t, name, projectParam.Name)
		assert.Equal(t, cwd, projectParam.Template.Root)
	})

	t.Run("panics with invalid data", func(t *testing.T) {
		for _, args := range [][]string{
			{"", ""},
			{"", "foo"},
			{"/foo/bar", ""},
		} {
			t.Run(fmt.Sprintf("for %s and %s", args[0], args[1]), func(t *testing.T) {
				// given
				svc := New(nil, nil, nil, nil)
				cwd := args[0]
				name := args[1]

				// expect
				defer func() {
					assert.NotNil(t, recover(), "The code did not panic")
				}()
				svc.CreateProject(cwd, name)
			})
		}
	})
}

func Test_SelectAndOpenProject(t *testing.T) {
	t.Run("runs selector and attaches to project when name is empty", func(t *testing.T) {
		// given
		projects := []*Project{{ID: "1234", Name: "foobar"}}
		projectNames := []string{projects[0].Name}

		slMock := new(MockSelector)
		slMock.On("SelectFrom", projectNames, mock.Anything).Return(projectNames[0], nil).Once()

		stMock := new(MockStorage)
		stMock.On("List").Return(projects, nil).Once()

		muMock := new(MockMultiplexer)
		muMock.On("AttachProject", projects[0]).Return(nil).Once()

		svc := New(slMock, muMock, stMock, nil)

		// when
		svc.SelectAndOpenProject("")

		// then
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
		muMock.AssertExpectations(t)
	})

	t.Run("tries to find project if name is provided", func(t *testing.T) {
		// given
		name := "foobar"
		project := &Project{ID: "1234", Name: name}

		stMock := new(MockStorage)
		stMock.On("Find", name).Return(project, nil).Once()

		muMock := new(MockMultiplexer)
		muMock.On("AttachProject", project).Return(nil).Once()

		svc := New(nil, muMock, stMock, nil)

		// when
		svc.SelectAndOpenProject(name)

		// then
		stMock.AssertExpectations(t)
		muMock.AssertExpectations(t)
	})

	t.Run("panics when project is not found", func(t *testing.T) {
		// given
		name := "foobar"
		err := ErrNotFound

		stMock := new(MockStorage)
		stMock.On("Find", name).Return(&Project{}, err).Once()

		svc := New(nil, nil, stMock, nil)

		// when
		defer func() {
			r := recover()
			assert.NotNil(t, r, "The code did not panic")
			assert.Equal(t, err, r, "The error should be %s was %s", err, r)
		}()
		svc.SelectAndOpenProject(name)

		// then
		stMock.AssertExpectations(t)
	})

	t.Run("exit gracefully when selector is cancelled", func(t *testing.T) {
		// given
		err := ErrSelectorCancelled
		listReturn := []*Project{{ID: "1234", Name: "foobar"}}

		slMock := new(MockSelector)
		slMock.On("SelectFrom", mock.Anything, mock.Anything).Return("", err).Once()

		stMock := new(MockStorage)
		stMock.On("List").Return(listReturn, nil).Once()

		svc := New(slMock, nil, stMock, nil)

		// when
		svc.SelectAndOpenProject("")

		// then
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
	})
}

func Test_DeleteProject(t *testing.T) {
	t.Run("runs selector and deletes project when name is empty", func(t *testing.T) {
		// given
		projects := []*Project{{ID: "1234", Name: "foobar"}}
		projectNames := []string{projects[0].Name}

		slMock := new(MockSelector)
		slMock.On("SelectFrom", projectNames, mock.Anything).Return(projectNames[0], nil).Once()

		stMock := new(MockStorage)
		stMock.On("List").Return(projects, nil).Once()
		stMock.On("Delete", projects[0].ID).Return(nil).Once()

		svc := New(slMock, nil, stMock, nil)

		// when
		svc.DeleteProject("")

		// then
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
	})

	t.Run("tries to find project if name is provided", func(t *testing.T) {
		// given
		name := "foobar"
		project := &Project{ID: "1234", Name: name}

		stMock := new(MockStorage)
		stMock.On("Find", name).Return(project, nil).Once()
		stMock.On("Delete", project.ID).Return(nil).Once()

		svc := New(nil, nil, stMock, nil)

		// when
		svc.DeleteProject(name)

		// then
		stMock.AssertExpectations(t)
	})

	t.Run("panics when project is not found", func(t *testing.T) {
		// given
		name := "foobar"
		err := ErrNotFound

		stMock := new(MockStorage)
		stMock.On("Find", name).Return(&Project{}, err).Once()

		svc := New(nil, nil, stMock, nil)

		// when
		defer func() {
			r := recover()
			assert.NotNil(t, r, "The code did not panic")
			assert.Equal(t, err, r, "The error should be %s was %s", err, r)
		}()
		svc.DeleteProject(name)

		// then
		stMock.AssertExpectations(t)
	})

	t.Run("exit gracefully when selector is cancelled", func(t *testing.T) {
		// given
		err := ErrSelectorCancelled
		listReturn := []*Project{{ID: "1234", Name: "foobar"}}

		slMock := new(MockSelector)
		slMock.On("SelectFrom", mock.Anything, mock.Anything).Return("", err).Once()

		stMock := new(MockStorage)
		stMock.On("List").Return(listReturn, nil).Once()

		svc := New(slMock, nil, stMock, nil)

		// when
		svc.DeleteProject("")

		// then
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
	})
}

func Test_EditProject(t *testing.T) {
	t.Run("runs selector and launches editor when name is empty", func(t *testing.T) {
		// given
		projects := []*Project{{ID: "1234", Name: "foobar"}}
		projectNames := []string{projects[0].Name}
		cwd := "/home/test"
		template := Template{Root: cwd, Windows: []Window{{Name: "main", Root: "/project"}}}
		projects[0].Template = &template
		templateFile := "/home/test/template.yaml"

		slMock := new(MockSelector)
		slMock.On("SelectFrom", projectNames, mock.Anything).Return(projectNames[0], nil).Once()

		stMock := new(MockStorage)
		stMock.On("List").Return(projects, nil).Once()
		stMock.On("PrepareTemplateFile", projects[0]).Return(templateFile, nil).Once()

		editorMock := new(MockEditorLauncher)
		editorMock.On("Open", templateFile).Return(nil).Once()

		svc := New(slMock, nil, stMock, editorMock)

		// when
		svc.EditProject("")

		// then
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
		editorMock.AssertExpectations(t)
	})

	t.Run("tries to find project if name is provided", func(t *testing.T) {
		// given
		name := "foobar"
		project := &Project{ID: "1234", Name: name}
		template := Template{Root: "/home/test", Windows: []Window{{Name: "main", Root: "/project"}}}
		project.Template = &template
		templateFile := "/home/test/template.yaml"

		stMock := new(MockStorage)
		stMock.On("Find", name).Return(project, nil).Once()
		stMock.On("PrepareTemplateFile", project).Return(templateFile, nil).Once()

		editorMock := new(MockEditorLauncher)
		editorMock.On("Open", templateFile).Return(nil).Once()

		svc := New(nil, nil, stMock, editorMock)

		// when
		svc.EditProject(name)

		// then
		stMock.AssertExpectations(t)
		editorMock.AssertExpectations(t)
	})

	t.Run("panics when project is not found", func(t *testing.T) {
		// given
		name := "foobar"
		err := ErrNotFound

		stMock := new(MockStorage)
		stMock.On("Find", name).Return(&Project{}, err).Once()

		svc := New(nil, nil, stMock, nil)

		// when
		defer func() {
			r := recover()
			assert.NotNil(t, r, "The code did not panic")
			assert.Equal(t, err, r, "The error should be %s was %s", err, r)
		}()
		svc.EditProject(name)

		// then
		stMock.AssertExpectations(t)
	})

	t.Run("exit gracefully when selector is cancelled", func(t *testing.T) {
		// given
		err := ErrSelectorCancelled
		listReturn := []*Project{{ID: "1234", Name: "foobar"}}

		slMock := new(MockSelector)
		slMock.On("SelectFrom", mock.Anything, mock.Anything).Return("", err).Once()

		stMock := new(MockStorage)
		stMock.On("List").Return(listReturn, nil).Once()

		svc := New(slMock, nil, stMock, nil)

		// when
		svc.EditProject("")

		// then
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
	})
}
