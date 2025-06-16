package multiplexer_test

import (
	"testing"
	"thop/internal/multiplexer"
	"thop/internal/problem"
	"thop/internal/types/command"
	"thop/internal/types/project"
	"thop/internal/types/template"
	"thop/internal/types/window"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTmuxClient struct {
	mock.Mock
}

func (m *MockTmuxClient) AttachSession(session multiplexer.SessionName) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockTmuxClient) SwitchSession(session multiplexer.SessionName) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockTmuxClient) HasSession(session multiplexer.SessionName) (bool, error) {
	args := m.Called(session)
	return args.Bool(0), args.Error(1)
}

func (m *MockTmuxClient) NewSession(
	session multiplexer.SessionName,
	root template.Root,
	windowName window.Name,
	windowRoot window.Root,
) error {
	args := m.Called(session, root, windowName, windowRoot)
	return args.Error(0)
}

func (m *MockTmuxClient) NewWindow(
	session multiplexer.SessionName,
	root template.Root,
	windowName window.Name,
	windowRoot window.Root,
) error {
	args := m.Called(session, root, windowName, windowRoot)
	return args.Error(0)
}

func (m *MockTmuxClient) SendKeys(
	session multiplexer.SessionName,
	windowName window.Name,
	keys command.Command,
) error {
	args := m.Called(session, windowName, keys)
	return args.Error(0)
}

func (m *MockTmuxClient) ListSessions() ([]multiplexer.SessionName, error) {
	args := m.Called()
	return args.Get(0).([]multiplexer.SessionName), args.Error(1)
}

func Test_AttachProject(t *testing.T) {
	t.Run("assembles and attaches to session if it doesn't exist", func(t *testing.T) {
		// given
		uuid := project.UUID("foo")
		name := project.Name("foo")
		sessionName := multiplexer.SessionName(name)
		root := template.Root("/home/test")
		window1Name := window.Name("main")
		window1Root := window.Root("/project")
		window2Name := window.Name("baz")
		window2Root := window.Root("")
		project := project.Project{
			UUID: uuid,
			Name: name,
			Template: template.Template{
				Root:     root,
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Name: window1Name,
						Root: "/project",
					},
					{
						Name:     window2Name,
						Commands: []command.Command{"ls"},
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", sessionName).Return(false, nil).Once()
		mockClient.On("NewSession", sessionName, root, window1Name, window1Root).Return(nil).Once()
		mockClient.On("NewWindow", sessionName, root, window2Name, window2Root).Return(nil).Once()
		mockClient.On("SendKeys", sessionName, window1Name, command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", sessionName, window2Name, command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", sessionName, window2Name, command.Command("ls")).Return(nil).Once()
		mockClient.On("AttachSession", sessionName).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("attaches if session exists", func(t *testing.T) {
		// given
		project := project.Project{
			UUID: "foo",
			Name: "foo",
			Template: template.Template{
				Root:     "/home/test",
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Name: "main",
						Root: "/project",
					},
					{
						Name:     "baz",
						Commands: []command.Command{"ls"},
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", multiplexer.SessionName("foo")).Return(true, nil).Once()
		mockClient.On("AttachSession", multiplexer.SessionName("foo")).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("assembles and switches to session if it doesn't exist and shell is in active session", func(t *testing.T) {
		// given
		uuid := project.UUID("foo")
		name := project.Name("foo")
		sessionName := multiplexer.SessionName(name)
		root := template.Root("/home/test")
		window1Name := window.Name("main")
		window1Root := window.Root("/project")
		window2Name := window.Name("baz")
		window2Root := window.Root("")
		project := project.Project{
			UUID: uuid,
			Name: name,
			Template: template.Template{
				Root:     root,
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Name: window1Name,
						Root: window1Root,
					},
					{
						Name:     window2Name,
						Commands: []command.Command{"ls"},
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", sessionName).Return(false, nil).Once()
		mockClient.On("NewSession", sessionName, root, window1Name, window1Root).Return(nil).Once()
		mockClient.On("NewWindow", sessionName, root, window2Name, window2Root).Return(nil).Once()
		mockClient.On("SendKeys", sessionName, window1Name, command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", sessionName, window2Name, command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", sessionName, window2Name, command.Command("ls")).Return(nil).Once()
		mockClient.On("SwitchSession", sessionName).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client:            mockClient,
			ActiveTmuxSession: "/home/test",
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("switches to session if it exist and shell is in active session", func(t *testing.T) {
		// given
		project := project.Project{
			UUID: "foo",
			Name: "foo",
			Template: template.Template{
				Root:     "/home/test",
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Name: "main",
						Root: "/project",
					},
					{
						Name:     "baz",
						Commands: []command.Command{"ls"},
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", multiplexer.SessionName("foo")).Return(true, nil).Once()
		mockClient.On("SwitchSession", multiplexer.SessionName("foo")).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client:            mockClient,
			ActiveTmuxSession: "/home/test",
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("uses template name for session if provided", func(t *testing.T) {
		// given
		project := project.Project{
			UUID: "foo",
			Name: "foo",
			Template: template.Template{
				Name: "bar",
				Root: "/home/test",
				Windows: []window.Window{
					{
						Name: "main",
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", multiplexer.SessionName("bar")).Return(false, nil).Once()
		mockClient.On("NewSession", multiplexer.SessionName("bar"), template.Root("/home/test"), window.Name("main"), window.Root("")).Return(nil).Once()
		mockClient.On("AttachSession", multiplexer.SessionName("bar")).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("returns error if project has no name", func(t *testing.T) {
		multiplexer := multiplexer.TmuxMultiplexer{
			Client: nil,
		}

		err := multiplexer.AttachProject(project.Project{Name: "", Template: template.Template{Name: ""}})
		assert.NotNil(t, err, "Expected error when project has no name")
	})
}

func Test_ListActiveSessions(t *testing.T) {
	t.Run("returns error if client fails to list sessions", func(t *testing.T) {
		// given
		mockClient := new(MockTmuxClient)
		mockClient.On("ListSessions").Return(
			[]multiplexer.SessionName{},
			multiplexer.ErrFailedToListSessions.WithMsg("foo"),
		).Once()

		m := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		_, err := m.ListActiveSessions()

		// then
		assert.Equal(t, multiplexer.ErrFailedToListSessions, err.(problem.Problem).Key)
		mockClient.AssertExpectations(t)
	})

	t.Run("returns empty list if client returns empty list", func(t *testing.T) {
		// given
		mockClient := new(MockTmuxClient)
		mockClient.On("ListSessions").Return([]multiplexer.SessionName{}, nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		sessions, err := multiplexer.ListActiveSessions()

		// then
		assert.Nil(t, err)
		assert.Equal(t, []project.Project(nil), sessions)
		mockClient.AssertExpectations(t)
	})

	t.Run("returns list of projects", func(t *testing.T) {
		// given
		mockClient := new(MockTmuxClient)
		mockClient.On("ListSessions").Return([]multiplexer.SessionName{"foo", "bar"}, nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		sessions, err := multiplexer.ListActiveSessions()

		// then
		assert.Nil(t, err)
		for i, session := range []project.Project{{Name: "foo"}, {Name: "bar"}} {
			assert.Equal(t, session, sessions[i])
		}
		mockClient.AssertExpectations(t)
	})
}
