package tmux_test

import (
	"testing"
	. "thop/dom/model"
	"thop/infra/tmux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTmuxClient struct {
	mock.Mock
}

func (m *MockTmuxClient) AttachSession(sessionName string) error {
	args := m.Called(sessionName)
	return args.Error(0)
}

func (m *MockTmuxClient) SwitchSession(sessionName string) error {
	args := m.Called(sessionName)
	return args.Error(0)
}

func (m *MockTmuxClient) HasSession(sessionName string) (bool, error) {
	args := m.Called(sessionName)
	return args.Bool(0), args.Error(1)
}

func (m *MockTmuxClient) NewSession(sessionName, sessionRoot, windowName, windowRoot string) error {
	args := m.Called(sessionName, sessionRoot, windowName, windowRoot)
	return args.Error(0)
}

func (m *MockTmuxClient) NewWindow(sessionName, sessionRoot, windowName, windowRoot string) error {
	args := m.Called(sessionName, sessionRoot, windowName, windowRoot)
	return args.Error(0)
}

func (m *MockTmuxClient) SendKeys(sessionName, windowName, keys string) error {
	args := m.Called(sessionName, windowName, keys)
	return args.Error(0)
}

func (m *MockTmuxClient) IsInTmuxSession() bool {
	args := m.Called()
	return args.Bool(0)
}

func Test_AttachProject(t *testing.T) {
	t.Run("assembles and attaches to session if it doesn't exist", func(t *testing.T) {
		// given
		project := Project{
			ID:   "foo",
			Name: "foo",
			Template: &Template{
				Root:     "/home/test",
				Commands: []string{"echo hello"},
				Windows: []Window{
					{
						Name: "main",
						Root: "/project",
					},
					{
						Name:     "baz",
						Commands: []string{"ls"},
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", "foo").Return(false, nil).Once()
		mockClient.On("NewSession", "foo", "/home/test", "main", "/project").Return(nil).Once()
		mockClient.On("NewWindow", "foo", "/home/test", "baz", "").Return(nil).Once()
		mockClient.On("SendKeys", "foo", "main", "echo hello").Return(nil).Once()
		mockClient.On("SendKeys", "foo", "baz", "echo hello").Return(nil).Once()
		mockClient.On("SendKeys", "foo", "baz", "ls").Return(nil).Once()
		mockClient.On("IsInTmuxSession").Return(false).Once()
		mockClient.On("AttachSession", "foo").Return(nil).Once()

		multiplexer := tmux.NewTmuxMultiplexer(mockClient)

		// when
		err := multiplexer.AttachProject(&project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("attaches if session exists", func(t *testing.T) {
		// given
		project := Project{
			ID:   "foo",
			Name: "foo",
			Template: &Template{
				Root:     "/home/test",
				Commands: []string{"echo hello"},
				Windows: []Window{
					{
						Name: "main",
						Root: "/project",
					},
					{
						Name:     "baz",
						Commands: []string{"ls"},
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", "foo").Return(true, nil).Once()
		mockClient.On("IsInTmuxSession").Return(false).Once()
		mockClient.On("AttachSession", "foo").Return(nil).Once()

		multiplexer := tmux.NewTmuxMultiplexer(mockClient)

		// when
		err := multiplexer.AttachProject(&project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("assembles and switches to session if it doesn't exist and shell is in active session", func(t *testing.T) {
		// given
		project := Project{
			ID:   "foo",
			Name: "foo",
			Template: &Template{
				Root:     "/home/test",
				Commands: []string{"echo hello"},
				Windows: []Window{
					{
						Name: "main",
						Root: "/project",
					},
					{
						Name:     "baz",
						Commands: []string{"ls"},
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", "foo").Return(false, nil).Once()
		mockClient.On("NewSession", "foo", "/home/test", "main", "/project").Return(nil).Once()
		mockClient.On("NewWindow", "foo", "/home/test", "baz", "").Return(nil).Once()
		mockClient.On("SendKeys", "foo", "main", "echo hello").Return(nil).Once()
		mockClient.On("SendKeys", "foo", "baz", "echo hello").Return(nil).Once()
		mockClient.On("SendKeys", "foo", "baz", "ls").Return(nil).Once()
		mockClient.On("IsInTmuxSession").Return(true).Once()
		mockClient.On("SwitchSession", "foo").Return(nil).Once()

		multiplexer := tmux.NewTmuxMultiplexer(mockClient)

		// when
		err := multiplexer.AttachProject(&project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("switches to session if it exist and shell is in active session", func(t *testing.T) {
		// given
		project := Project{
			ID:   "foo",
			Name: "foo",
			Template: &Template{
				Root:     "/home/test",
				Commands: []string{"echo hello"},
				Windows: []Window{
					{
						Name: "main",
						Root: "/project",
					},
					{
						Name:     "baz",
						Commands: []string{"ls"},
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", "foo").Return(true, nil).Once()
		mockClient.On("IsInTmuxSession").Return(true).Once()
		mockClient.On("SwitchSession", "foo").Return(nil).Once()

		multiplexer := tmux.NewTmuxMultiplexer(mockClient)

		// when
		err := multiplexer.AttachProject(&project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("uses template name for session if provided", func(t *testing.T) {
		// given
		project := Project{
			ID:   "foo",
			Name: "foo",
			Template: &Template{
				Name: "bar",
				Root: "/home/test",
				Windows: []Window{
					{
						Name: "main",
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", "bar").Return(false, nil).Once()
		mockClient.On("NewSession", "bar", "/home/test", "main", "").Return(nil).Once()
		mockClient.On("IsInTmuxSession").Return(false).Once()
		mockClient.On("AttachSession", "bar").Return(nil).Once()

		multiplexer := tmux.NewTmuxMultiplexer(mockClient)

		// when
		err := multiplexer.AttachProject(&project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("returns error if project has no name", func(t *testing.T) {
		multiplexer := tmux.NewTmuxMultiplexer(nil)

		err := multiplexer.AttachProject(&Project{Name: "", Template: &Template{Name: ""}})
		assert.NotNil(t, err, "Expected error when project has no name")
	})
}
