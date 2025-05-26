package cli_test

import (
	"phopper/src/app/cli"
	. "phopper/test/utils"
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

func Test_CLI(t *testing.T) {
	t.Run("select command", func(t *testing.T) {
		t.Run("select project", func(t *testing.T) {
			variants := [][]string{
				{"select"},
				{"s"},
				{},
			}

			for _, args := range variants {
				// given
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				Assert(t, svc.SelectAndOpenProjectCalls == 1, "SelectAndOpenProject should be called once")
				Assert(t, svc.SelectAndOpenProjectParam1 == "", "Name should be empty")
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
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				Assert(t, svc.SelectAndOpenProjectCalls == 1, "SelectAndOpenProject should be called once")
				param := svc.SelectAndOpenProjectParam1
				Assert(t, param == name, "Name should be %s is %s", name, param)
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
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				Assert(t, svc.CreateProjectCalls == 1, "CreateProject should be called once")
				param1 := svc.CreateProjectParam1
				param2 := svc.CreateProjectParam2
				Assert(t, param1 == param2, "Name should default to current working directory")
			}
		})

		t.Run("create project with name", func(t *testing.T) {
			name := "foobar"

			for _, args := range variants {
				// given
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(append(args, name))

				// then
				Assert(t, svc.CreateProjectCalls == 1, "CreateProject should be called once")
				param2 := svc.CreateProjectParam2
				Assert(t, param2 == name, "Name should be %s is %s", name, param2)
			}
		})

		t.Run("create project with name and cwd", func(t *testing.T) {
			name := "foobar"
			cwd := "/home/test"

			for _, args := range variants {
				// given
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(append(args, name, cwd))

				// then
				Assert(t, svc.CreateProjectCalls == 1, "CreateProject should be called once")

				param1 := svc.CreateProjectParam1
				Assert(t, param1 == cwd, "Cwd should be %s is %s", cwd, param1)
				param2 := svc.CreateProjectParam2
				Assert(t, param2 == name, "Name should be %s is %s", name, param2)
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
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				Assert(t, svc.DeleteProjectCalls == 1, "DeleteProject should be called once")
				param := svc.DeleteProjectParam1
				Assert(t, param == "", "Name should be empty")
			}
		})

		t.Run("delete project with name", func(t *testing.T) {
			name := "foobar"

			for _, args := range variants {
				// given
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(append(args, name))

				// then
				Assert(t, svc.DeleteProjectCalls == 1, "DeleteProject should be called once")
				param := svc.DeleteProjectParam1
				Assert(t, param == name, "Name should be %s is %s", name, param)
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
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(args)

				// then
				Assert(t, svc.EditProjectCalls == 1, "EditProject should be called once")
				Assert(t, svc.EditProjectParam1 == "", "Name should be empty")
			}
		})

		t.Run("delete project with name", func(t *testing.T) {
			name := "foobar"

			for _, args := range variants {
				// given
				svc := &MockService{}
				cli := cli.NewCli(svc)

				// when
				cli.Run(append(args, name))

				// then
				Assert(t, svc.EditProjectCalls == 1, "EditProject should be called once")
				param := svc.EditProjectParam1
				Assert(t, param == name, "Name should be %s is %s", name, param)
			}
		})
	})
}
