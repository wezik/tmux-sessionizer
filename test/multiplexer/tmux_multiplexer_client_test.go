package multiplexer_test

import (
	"errors"
	"os/exec"
	"testing"
	"thop/internal/multiplexer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommandExecutor struct {
	mock.Mock
	ExecutedCommands [][]string
}

func (m *MockCommandExecutor) Execute(cmd *exec.Cmd) (string, int, error) {
	args := m.Called(cmd)
	m.ExecutedCommands = append(m.ExecutedCommands, cmd.Args)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockCommandExecutor) ExecuteInteractive(cmd *exec.Cmd) (int, error) {
	args := m.Called(cmd)
	m.ExecutedCommands = append(m.ExecutedCommands, cmd.Args)
	return args.Int(0), args.Error(1)
}

func Test_Client_AttachSession(t *testing.T) {
	t.Run("returns error if session name is empty", func(t *testing.T) {
		// given
		client := &multiplexer.TmuxClientImpl{}

		// when
		err := client.AttachSession("")

		// then
		assert.True(t, multiplexer.ErrInvalidTemplateArgs.Equal(err))
	})

	t.Run("attaches to session", func(t *testing.T) {
		// given
		mockExecutor := new(MockCommandExecutor)

		mockExecutor.On("Execute", mock.Anything).Return("", 0, nil)

		expectedCmd := [][]string{
			{"tmux", "attach", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: mockExecutor,
		}

		// when
		err := client.AttachSession("mysession")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, mockExecutor.ExecutedCommands)
	})
}

func Test_Client_SwitchSession(t *testing.T) {
	t.Run("returns error if session name is empty", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// when
		err := client.SwitchSession("")

		// then
		assert.True(t, multiplexer.ErrInvalidTemplateArgs.Equal(err))
	})

	t.Run("switches session", func(t *testing.T) {
		// given
		mockExecutor := new(MockCommandExecutor)
		mockExecutor.On("Execute", mock.Anything).Return("", 0, nil)

		expectedCmd := [][]string{
			{"tmux", "switch", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: mockExecutor,
		}

		// when
		err := client.SwitchSession("mysession")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, mockExecutor.ExecutedCommands)
	})

}

func Test_Client_HasSession(t *testing.T) {
	t.Run("returns error if session name is empty", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// when
		_, err := client.HasSession("")

		// then
		assert.True(t, multiplexer.ErrInvalidTemplateArgs.Equal(err))
	})

	t.Run("returns false on exit code 1 gracefully", func(t *testing.T) {
		// given
		mockExecutor := new(MockCommandExecutor)
		mockExecutor.On("Execute", mock.Anything).Return("", 1, errors.New("exit code 1"))
		expectedCmd := [][]string{
			{"tmux", "has-session", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: mockExecutor,
		}

		// when
		exists, err := client.HasSession("mysession")

		// then
		assert.Nil(t, err)
		assert.False(t, exists)
		assert.Equal(t, expectedCmd, mockExecutor.ExecutedCommands)
	})

	t.Run("returns true if session exists", func(t *testing.T) {
		// given
		mockExecutor := new(MockCommandExecutor)
		mockExecutor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{"tmux", "has-session", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: mockExecutor,
		}

		// when
		exists, err := client.HasSession("mysession")

		// then
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, expectedCmd, mockExecutor.ExecutedCommands)
	})
}

func Test_Client_NewSession(t *testing.T) {
	t.Run("returns error when missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// expect
		err := client.NewSession("", "root", "win", "")
		assert.NotNil(t, err, "expected error for empty session name")

		// and
		err = client.NewSession("sess", "", "win", "")
		assert.NotNil(t, err, "expected error for empty session root")

		// and
		err = client.NewSession("sess", "root", "", "")
		assert.NotNil(t, err, "expected error for empty window name")
	})

	t.Run("returns error if session already exists", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 1, errors.New("exit code 1"))
		expectedCmd := [][]string{
			{
				"tmux",
				"new-session",
				"-d",
				"-s",
				"mysession",
				"-c",
				"/home/test",
				"-n",
				"main",
				"cd /project && exec $SHELL",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.NewSession("mysession", "/home/test", "main", "/project")

		// then
		assert.NotNil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new session", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-session",
				"-d",
				"-s",
				"mysession",
				"-c",
				"/home/test",
				"-n",
				"main",
				"cd /project && exec $SHELL",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.NewSession("mysession", "/home/test", "main", "/project")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("defaults main window root to session root if empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-session",
				"-d",
				"-s",
				"mysession",
				"-c",
				"/home/test",
				"-n",
				"main",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.NewSession("mysession", "/home/test", "main", "")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_TmuxClient_SendKeys(t *testing.T) {
	t.Run("returns error when missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// expect
		err := client.SendKeys("", "win", "ls")
		assert.NotNil(t, err, "expected error for empty session name")

		// and
		err = client.SendKeys("sess", "", "ls")
		assert.NotNil(t, err, "expected error for empty window name")

		// and
		err = client.SendKeys("sess", "win", "")
		assert.NotNil(t, err, "expected error for empty keys")
	})

	t.Run("sends keys to window", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"send-keys",
				"-t",
				"mysession:main",
				"ls",
				"C-m",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SendKeys("mysession", "main", "ls")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_Client_NewWindow(t *testing.T) {
	t.Run("returns error when missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// expect
		err := client.NewWindow("", "/project", "window", "/root")
		assert.NotNil(t, err, "expected error for empty session name")

		// and
		err = client.NewWindow("sess", "/project", "", "/root")
		assert.NotNil(t, err, "expected error for empty window name")

		// and
		err = client.NewWindow("sess", "", "window", "/root")
		assert.NotNil(t, err, "expected error for empty session root")
	})

	t.Run("includes default seession root if window root is empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-window",
				"-d",
				"-t",
				"mysession",
				"-n",
				"main",
				"-c",
				"/home/test",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.NewWindow("mysession", "/home/test", "main", "")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new window", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-window",
				"-d",
				"-t",
				"mysession",
				"-n",
				"main",
				"-c",
				"/project",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.NewWindow("mysession", "/home/test", "main", "/project")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_Client_ListSessions(t *testing.T) {
	t.Run("returns list of sessions", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("foo\nbar\nbaz\n", 0, nil).Once()

		expectedCmd := [][]string{
			{"tmux", "list-sessions", "-F", "#S"},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		sessions, err := client.ListSessions()

		// then
		assert.Nil(t, err)
		assert.Equal(t, []multiplexer.SessionName{"foo", "bar", "baz"}, sessions)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("returns mapped error if command fails", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 1, errors.New("exit code 1")).Once()

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		_, err := client.ListSessions()

		// then
		assert.True(t, multiplexer.ErrFailedToListSessions.Equal(err))
		executor.AssertExpectations(t)
	})
}
