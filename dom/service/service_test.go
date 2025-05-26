package service_test

import (
	"fmt"
	. "phopper/dom/model"
	. "phopper/dom/service"
	. "phopper/dom/utils"
	"slices"
	"testing"
)

type MockSelector struct {
	SelectFromParam1 []string
	SelectFromParam2 string
	SelectFromCalls  int
	SelectFromReturn string
	SelectFromErr    error
}

func (s *MockSelector) SelectFrom(items []string, prompt string) (string, error) {
	s.SelectFromParam1 = items
	s.SelectFromParam2 = prompt
	s.SelectFromCalls++
	return s.SelectFromReturn, s.SelectFromErr
}

type MockMultiplexer struct {
	AttachProjectParam1 *Project
	AttachProjectCalls  int
}

func (s *MockMultiplexer) AttachProject(p *Project) error {
	s.AttachProjectParam1 = p
	s.AttachProjectCalls++
	return nil
}

type MockStorage struct {
	ListCalls  int
	ListReturn []*Project
	ListErr    error

	FindParam1 string
	FindCalls  int
	FindReturn *Project
	FindErr    error

	SaveParam1 *Project
	SaveCalls  int

	DeleteParam1 string
	DeleteCalls  int
}

func (s *MockStorage) List() ([]*Project, error) {
	s.ListCalls++
	return s.ListReturn, s.ListErr
}

func (s *MockStorage) Find(name string) (*Project, error) {
	s.FindParam1 = name
	s.FindCalls++
	return s.FindReturn, s.FindErr
}

func (s *MockStorage) Save(t *Project) error {
	s.SaveParam1 = t
	s.SaveCalls++
	return nil
}

func (s *MockStorage) Delete(uuid string) error {
	s.DeleteParam1 = uuid
	s.DeleteCalls++
	return nil
}

func Test_Service(t *testing.T) {
	t.Run("create project", func(t *testing.T) {
		t.Run("creates project", func(t *testing.T) {
			// given
			st := &MockStorage{}
			svc := NewService(nil, nil, st)
			cwd := "/home/test"
			name := "foobar"

			// when
			svc.CreateProject(cwd, name)

			// then
			Assert(t, st.SaveCalls == 1, "Save should be called once")

			param1 := st.SaveParam1
			Assert(t, param1.Name == name, "Saved project should have name %s has %s", name, param1.Name)

			template := param1.Template
			Assert(t, template.Root == cwd, "Root should be %s is %s", cwd, template.Root)
			Assert(t, template.Name != name, "Name should be %s is %s", name, template.Name)
			Assert(t, len(template.Windows) == 1, "Saved project should have one window")
		})

		t.Run("panics with invalid data", func(t *testing.T) {
			for _, args := range [][]string{
				{"", ""},
				{"", "foo"},
				{"/foo/bar", ""},
			} {
				t.Run(fmt.Sprintf("for %s and %s", args[0], args[1]), func(t *testing.T) {
					// given
					svc := NewService(nil, nil, nil)
					cwd := args[0]
					name := args[1]

					// expect
					defer func() {
						Assert(t, recover() != nil, "The code did not panic")
					}()
					svc.CreateProject(cwd, name)
				})
			}
		})
	})

	t.Run("select and open project", func(t *testing.T) {
		t.Run("finds and opens project with multiplexer", func(t *testing.T) {
			// given
			name := "foobar"
			project := &Project{ID: "1234", Name: name}

			sl := &MockSelector{}

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.FindReturn = project
			svc := NewService(sl, mu, st)

			// when
			svc.SelectAndOpenProject(name)

			// then
			Assert(t, st.FindCalls == 1, "Find should be called once")

			param1 := st.FindParam1
			Assert(t, param1 == name, "Find param name should be %s is %s", name, param1)

			Assert(t, sl.SelectFromCalls == 0, "Selector should not be called")

			Assert(t, mu.AttachProjectCalls == 1, "The project should be attached")
			paramAttach := mu.AttachProjectParam1
			Assert(t, paramAttach == project, "Attach param name should be %s is %s", project.Name, paramAttach.Name)
		})

		t.Run("selects from selector and opens project with multiplexer", func(t *testing.T) {
			// given
			name := "foobar"
			project := &Project{ID: "1234", Name: name}
			projects := []*Project{project}
			projectNames := []string{project.Name}

			sl := &MockSelector{}
			sl.SelectFromReturn = project.Name

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.ListReturn = projects

			svc := NewService(sl, mu, st)

			// when
			svc.SelectAndOpenProject("")

			// then
			Assert(t, st.FindCalls == 0, "Find should not be called")

			Assert(t, sl.SelectFromCalls == 1, "Selector should be called once")

			param1 := sl.SelectFromParam1
			slicesEqual := slices.Equal(projectNames, param1)
			Assert(t, slicesEqual, "Selector items param should be %s is %s", projectNames, param1)

			attachParam := mu.AttachProjectParam1
			Assert(t, mu.AttachProjectCalls == 1, "The project should be attached")
			Assert(t, attachParam == project, "attach project should be %s is %s", project.ID, attachParam.ID)
		})

		t.Run("panics when project is not found", func(t *testing.T) {
			// given
			name := "foobar"
			err := ErrNotFound

			sl := &MockSelector{}

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.FindErr = err

			svc := NewService(sl, mu, st)

			// when
			defer func() {
				r := recover()
				Assert(t, r != nil, "The code did not panic")
				Assert(t, r.(error) == ErrNotFound, "The error should be %s was %s", err, r)
			}()
			svc.SelectAndOpenProject(name)

			// then
			Assert(t, st.FindCalls == 1, "Find should be called once")
			findParam := st.FindParam1
			Assert(t, findParam == name, "Find param name should be %s is %s", name, findParam)

			Assert(t, sl.SelectFromCalls == 0, "Selector should not be called")

			Assert(t, mu.AttachProjectCalls == 0, "The project should not be attached")
		})

		t.Run("exit gracefully when selector is cancelled", func(t *testing.T) {
			// given
			err := ErrSelectorCancelled
			listReturn := []*Project{{ID: "1234", Name: "foobar"}}

			sl := &MockSelector{}
			sl.SelectFromErr = err

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.ListReturn = listReturn

			svc := NewService(sl, mu, st)

			// when
			svc.SelectAndOpenProject("")

			// then
			Assert(t, st.FindCalls == 0, "Find should not be called")
			Assert(t, sl.SelectFromCalls == 1, "Selector should be called once")
			Assert(t, mu.AttachProjectCalls == 0, "The project should not be attached")
		})
	})

	t.Run("delete project", func(t *testing.T) {
		t.Run("finds and deletes project", func(t *testing.T) {
			// given
			name := "foobar"
			project := &Project{ID: "1234", Name: name}

			sl := &MockSelector{}

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.FindReturn = project

			svc := NewService(sl, mu, st)

			// when
			svc.DeleteProject(name)

			// then
			Assert(t, st.FindCalls == 1, "Find should be called once")

			findParam := st.FindParam1
			Assert(t, findParam == name, "Find param name should be %s is %s", name, findParam)

			Assert(t, sl.SelectFromCalls == 0, "Selector should not be called")

			Assert(t, st.DeleteCalls == 1, "Delete should be called once")

			deleteParam := st.DeleteParam1
			Assert(t, deleteParam == project.ID, "Delete param ID should be %s is %s", project.ID, deleteParam)
		})

		t.Run("selects from selector and deletes project", func(t *testing.T) {
			// given
			name := "foobar"
			project := &Project{ID: "1234", Name: name}
			projects := []*Project{project}
			projectNames := []string{project.Name}

			sl := &MockSelector{}
			sl.SelectFromReturn = project.Name

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.ListReturn = projects

			svc := NewService(sl, mu, st)

			// when
			svc.DeleteProject("")

			// then
			Assert(t, st.FindCalls == 0, "Find should not be called")

			Assert(t, sl.SelectFromCalls == 1, "Selector should be called once")

			selectFromParam := sl.SelectFromParam1
			slicesEqual := slices.Equal(projectNames, selectFromParam)
			Assert(t, slicesEqual, "Selector items param should be %s is %s", projectNames, selectFromParam)

			deleteParam := st.DeleteParam1
			Assert(t, deleteParam == project.ID, "Delete param ID should be %s is %s", project.ID, deleteParam)
		})

		t.Run("panics when project is not found", func(t *testing.T) {
			// given
			name := "foobar"
			err := ErrNotFound

			sl := &MockSelector{}

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.FindErr = err

			svc := NewService(sl, mu, st)

			// when
			defer func() {
				r := recover()
				Assert(t, r != nil, "The code did not panic")
				Assert(t, r.(error) == ErrNotFound, "The error should be %s was %s", err, r)
			}()
			svc.DeleteProject(name)

			// then
			Assert(t, st.FindCalls == 1, "Find should be called once")
			findParam := st.FindParam1
			Assert(t, findParam == name, "Find param name shoud be %s is %s", name, findParam)

			Assert(t, sl.SelectFromCalls == 0, "Selector should not be called")
			Assert(t, st.DeleteCalls == 0, "Delete should not be called")
		})

		t.Run("exit gracefully when selector is cancelled", func(t *testing.T) {
			// given
			err := ErrSelectorCancelled
			listReturn := []*Project{{ID: "1234", Name: "foobar"}}

			sl := &MockSelector{}
			sl.SelectFromErr = err

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.ListReturn = listReturn

			svc := NewService(sl, mu, st)

			// when
			svc.DeleteProject("")

			// then
			Assert(t, st.FindCalls == 0, "Find should not be called")
			Assert(t, sl.SelectFromCalls == 1, "Selector should be called once")
			Assert(t, st.DeleteCalls == 0, "Delete should not be called")
		})
	})
}
