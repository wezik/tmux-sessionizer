package cli_test

import (
	"phopper/src/app/cli"
	"phopper/src/domain/model"
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

func (s *MockService) CreateProject(cwd, name string) *model.Project {
	s.CreateProjectParam1 = cwd
	s.CreateProjectParam2 = name
	s.CreateProjectCalls++
	return nil
}

func TestSelectCmd(t *testing.T) {
	t.Parallel()

	t.Run("select project without name", func(t *testing.T) {
		t.Parallel()

		// given
		svc := mockService()
		cli := cli.NewCli(svc)
		args := []string{"select"}

		// when
		cli.Run(args)

		// then
		if svc.SelectAndOpenProjectCalls != 1 {
			t.Errorf("SelectAndOpenProject should be called once")
		}

		if svc.SelectAndOpenProjectParam1 != "" {
			t.Errorf("SelectAndOpenProject should be called with empty string")
		}
	})

	t.Run("select project with name", func(t *testing.T) {
		t.Parallel()

		// given
		svc := mockService()
		cli := cli.NewCli(svc)
		args := []string{"select", "test-name"}

		// when
		cli.Run(args)

		// then
		if svc.SelectAndOpenProjectCalls != 1 {
			t.Errorf("SelectAndOpenProject should be called once")
		}

		if svc.SelectAndOpenProjectParam1 != "test-name" {
			t.Errorf("SelectAndOpenProject should be called with test-name")
		}
	})
}

func TestCreateCmd(t *testing.T) {
	t.Parallel()

	t.Run("create project without name", func(t *testing.T) {
		t.Parallel()

		// given
		svc := mockService()
		cli := cli.NewCli(svc)
		args := []string{"create"}

		// when
		cli.Run(args)

		// then
		if svc.CreateProjectCalls != 1 {
			t.Errorf("CreateProject should be called once")
		}
		if svc.CreateProjectParam2 != svc.CreateProjectParam1 {
			t.Errorf("Name should default to current working directory")
		}
	})

	t.Run("create project with name", func(t *testing.T) {
		t.Parallel()

		// given
		svc := mockService()
		cli := cli.NewCli(svc)
		args := []string{"create", "test-name"}

		// when
		cli.Run(args)

		// then
		if svc.CreateProjectCalls != 1 {
			t.Errorf("CreateProject should be called once")
		}
		if svc.CreateProjectParam2 != "test-name" {
			t.Errorf("CreateProject should be called with test-name")
		}
	})

	t.Run("create project wih name and cwd", func(t *testing.T) {
		t.Parallel()
		// given
		svc := mockService()
		cli := cli.NewCli(svc)
		args := []string{"create", "test-name", "/home/test"}

		// when
		cli.Run(args)

		// then
		if svc.CreateProjectCalls != 1 {
			t.Errorf("CreateProject should be called once")
		}

		if svc.CreateProjectParam1 != "/home/test" {
			t.Errorf("CreateProject should be called with /home/test")
		}

		if svc.CreateProjectParam2 != "test-name" {
			t.Errorf("CreateProject should be called with test-name")
		}
	})
}
