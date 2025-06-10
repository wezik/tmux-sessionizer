package cli_test

import (
	"testing"
	. "thop/cli"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) SelectAndOpenProject(name string) {
	s.Called(name)
}

func (s *MockService) CreateProject(cwd, name string) {
	s.Called(cwd, name)
}

func (s *MockService) DeleteProject(name string) {
	s.Called(name)
}

func (s *MockService) EditProject(name string) {
	s.Called(name)
}

func Test_Select_Command(t *testing.T) {
	t.Run("selects project with empty name", func(t *testing.T) {
		variants := [][]string{
			{"select"},
			{"s"},
			{},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("SelectAndOpenProject", "").Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})

	t.Run("selects project with name", func(t *testing.T) {
		variants := [][]string{
			{"select", "foobar"},
			{"s", "foobar"},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("SelectAndOpenProject", "foobar").Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})
}

func Test_Create_Command(t *testing.T) {
	t.Run("creates project with no args", func(t *testing.T) {
		variants := [][]string{
			{"create"},
			{"c"},
			{"a"},
			{"add"},
			{"append"},
			{"new"},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("CreateProject", mock.Anything, mock.Anything).Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})

	t.Run("creates project with name", func(t *testing.T) {
		variants := [][]string{
			{"create", "foobar"},
			{"c", "foobar"},
			{"a", "foobar"},
			{"add", "foobar"},
			{"append", "foobar"},
			{"new", "foobar"},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("CreateProject", mock.Anything, "foobar").Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})

	t.Run("creates project with name and cwd", func(t *testing.T) {
		variants := [][]string{
			{"create", "foobar", "/home/test"},
			{"c", "foobar", "/home/test"},
			{"a", "foobar", "/home/test"},
			{"add", "foobar", "/home/test"},
			{"append", "foobar", "/home/test"},
			{"new", "foobar", "/home/test"},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("CreateProject", "/home/test", "foobar").Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})
}

func Test_Delete_Command(t *testing.T) {
	t.Run("deletes project with no args", func(t *testing.T) {
		variants := [][]string{
			{"delete"},
			{"d"},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("DeleteProject", mock.Anything).Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})

	t.Run("deletes project with name", func(t *testing.T) {
		variants := [][]string{
			{"delete", "foobar"},
			{"d", "foobar"},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("DeleteProject", "foobar").Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})
}

func Test_Edit_Command(t *testing.T) {
	t.Run("edits project with no args", func(t *testing.T) {
		variants := [][]string{
			{"edit"},
			{"e"},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("EditProject", mock.Anything).Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})

	t.Run("edits project with name", func(t *testing.T) {
		variants := [][]string{
			{"edit", "foobar"},
			{"e", "foobar"},
		}

		for _, args := range variants {
			// given
			svcMock := new(MockService)
			svcMock.On("EditProject", "foobar").Once()

			cli := NewCli(svcMock)

			// when
			cli.Run(args)

			// then
			svcMock.AssertExpectations(t)
		}
	})
}
