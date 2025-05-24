package cli_test

import (
	"phopper/src/app/cli"
	"testing"
)

type MockService struct {
	SelectAndOpenProjectParam1 string
	SelectAndOpenProjectCalls  int

	CreateProjectParam1 string
	CreateProjectParam2 string
	CreateProjectCalls  int

	DeleteProjectParam1 string
	DeleteProjectCalls  int

	EditProjectParam1 string
	EditProjectParam2 string
	EditProjectCalls  int
}

func mockService() *MockService {
	return &MockService{}
}

func (s *MockService) SelectAndOpenProject(name string) {
	s.SelectAndOpenProjectParam1 = name
	s.SelectAndOpenProjectCalls++
}

func (s *MockService) CreateProject(cwd, name string) {
	s.CreateProjectParam1 = cwd
	s.CreateProjectParam2 = name
	s.CreateProjectCalls++
}

func (s *MockService) DeleteProject(name string) {
	s.DeleteProjectParam1 = name
	s.DeleteProjectCalls++
}

func (s *MockService) EditProject(name string) {
	s.EditProjectParam1 = name
	s.EditProjectCalls++
}

func TestRun(t *testing.T) {
	t.Run("select command", func(t *testing.T) {
		t.Run("select project", func(t *testing.T) {
			variants := [][]string{
				{"select"},
				{"s"},
				{},
			}

			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				if svc.SelectAndOpenProjectCalls != 1 {
					t.Errorf("SelectAndOpenProject should be called once")
				}
				if svc.SelectAndOpenProjectParam1 != "" {
					t.Errorf("SelectAndOpenProject should be called with empty string")
				}
			}
		})

		t.Run("select project with name", func(t *testing.T) {
			name := "foobar"
			variants := [][]string{
				{"select", name},
				{"s", name},
			}

			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				if svc.SelectAndOpenProjectCalls != 1 {
					t.Errorf("SelectAndOpenProject should be called once")
				}

				if svc.SelectAndOpenProjectParam1 != name {
					t.Errorf("SelectAndOpenProject should be called with %s", name)
				}
			}
		})
	})

	t.Run("create command", func(t *testing.T) {
		variants := [][]string{
			{"create"},
			{"c"},
			{"a"},
			{"add"},
			{"append"},
			{"new"},
		}

		t.Run("create project", func(t *testing.T) {
			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				if svc.CreateProjectCalls != 1 {
					t.Errorf("CreateProject should be called once")
				}

				if svc.CreateProjectParam2 != svc.CreateProjectParam1 {
					t.Errorf("Name should default to current working directory")
				}
			}
		})

		t.Run("create project with name", func(t *testing.T) {
			name := "foobar"

			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(append(args, name))

				// then
				if svc.CreateProjectCalls != 1 {
					t.Errorf("CreateProject should be called once")
				}

				if svc.CreateProjectParam2 != name {
					t.Errorf("Name should be %s", name)
				}
			}
		})

		t.Run("create project with name and cwd", func(t *testing.T) {
			name := "foobar"
			cwd := "/home/test"

			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(append(args, name, cwd))

				// then
				if svc.CreateProjectCalls != 1 {
					t.Errorf("CreateProject should be called once")
				}

				if svc.CreateProjectParam1 != cwd {
					t.Errorf("Cwd should be %s", cwd)
				}

				if svc.CreateProjectParam2 != name {
					t.Errorf("Name should be %s", name)
				}
			}
		})
	})

	t.Run("delete command", func(t *testing.T) {
		variants := [][]string{
			{"delete"},
			{"d"},
		}

		t.Run("delete project", func(t *testing.T) {
			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				if svc.DeleteProjectCalls != 1 {
					t.Errorf("DeleteProject should be called once")
				}

				if svc.DeleteProjectParam1 != "" {
					t.Errorf("DeleteProject should be called with empty string")
				}
			}
		})

		t.Run("delete project with name", func(t *testing.T) {
			name := "foobar"

			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(append(args, name))

				// then
				if svc.DeleteProjectCalls != 1 {
					t.Errorf("DeleteProject should be called once")
				}

				if svc.DeleteProjectParam1 != name {
					t.Errorf("DeleteProject should be called with %s", name)
				}
			}
		})
	})

	t.Run("edit command", func(t *testing.T) {
		variants := [][]string{
			{"edit"},
			{"e"},
		}

		t.Run("edit project", func(t *testing.T) {
			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				if svc.EditProjectCalls != 1 {
					t.Errorf("EditProject should be called once")
				}

				if svc.EditProjectParam1 != "" {
					t.Errorf("EditProject should be called with empty string")
				}
			}
		})

		t.Run("delete project with name", func(t *testing.T) {
			name := "foobar"

			for _, args := range variants {
				// given
				svc := mockService()
				cli := cli.NewCli(svc)

				// when
				cli.Run(append(args, name))

				// then
				if svc.EditProjectCalls != 1 {
					t.Errorf("EditProject should be called once")
				}

				if svc.EditProjectParam1 != name {
					t.Errorf("EditProject should be called with %s", name)
				}
			}
		})
	})
}
