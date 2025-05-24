package cli_test

import (
	"phopper/src/app/cli"
	"phopper/src/domain/model"
	"testing"
)

func newCli (svc *MockService) *cli.Cli {
	return cli.NewCli(
		svc,
	)
}

type MockService struct {
	CreateProjectCalls int
	CreateProjectParam1 string
	CreateProjectParam2 string
}

func mockService() *MockService {
	return &MockService{}
}

func (s *MockService) CreateProject(cwd string, name string) (*model.Project, error) {
	s.CreateProjectCalls++
	s.CreateProjectParam1 = cwd
	s.CreateProjectParam2 = name
	return nil, nil
}

func TestCreateCmd(t *testing.T) {
	t.Parallel()

	t.Run("create project without name", func(t *testing.T) {
		t.Parallel()
		// given
		svc := mockService()
		cli := newCli(svc)
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
		cli := newCli(svc)
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
		cli := newCli(svc)
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
