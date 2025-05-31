package tmux_test

import (
	"errors"
	"os/exec"
	"testing"

	. "thop/dom/model"
	. "thop/dom/service"
	. "thop/infra/tmux"
)

type MockCommandExecutor struct {
	ExecutedCommands []*exec.Cmd
	ReturnStdout     string
	ReturnStderr     string
	ReturnErr        error
	ReturnExitCode   int
}

func (m *MockCommandExecutor) Execute(cmd *exec.Cmd) (string, int, error) {
	m.ExecutedCommands = append(m.ExecutedCommands, cmd)
	return m.ReturnStdout, m.ReturnExitCode, m.ReturnErr
}

func (m *MockCommandExecutor) ExecuteInteractive(cmd *exec.Cmd) (int, error) {
	m.ExecutedCommands = append(m.ExecutedCommands, cmd)
	return m.ReturnExitCode, m.ReturnErr
}

func assert(t *testing.T, condition bool, msg string, args ...any) {
	if !condition {
		t.Fatalf(msg, args...)
	}
}

func Test_TmuxClient(t *testing.T) {
	t.Run("AttachSession returns error if session name is empty", func(t *testing.T) {
		client := NewTmuxClient()
		err := client.AttachSession(&MockCommandExecutor{}, "")
		assert(t, err != nil, "Expected error when session name is empty")
	})

	t.Run("AttachSession calls correct tmux command", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := NewTmuxClient()

		err := client.AttachSession(executor, "mysession")

		assert(t, err == nil, "Expected no error")
		assert(t, len(executor.ExecutedCommands) == 1, "Expected 1 command execution")
		assert(t, executor.ExecutedCommands[0].Args[0] == "tmux", "Expected tmux command")
		assert(t, executor.ExecutedCommands[0].Args[1] == "attach", "Expected attach")
		assert(t, executor.ExecutedCommands[0].Args[3] == "mysession", "Expected session name")
	})

	t.Run("SwitchSession returns error if session name is empty", func(t *testing.T) {
		client := NewTmuxClient()
		executor := &MockCommandExecutor{}

		err := client.SwitchSession(executor, "")
		assert(t, err != nil, "Expected error when session name is empty")
	})

	t.Run("SwitchSession calls correct tmux command", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := NewTmuxClient()

		err := client.SwitchSession(executor, "mysession")

		assert(t, err == nil, "Expected no error")
		assert(t, len(executor.ExecutedCommands) == 1, "Expected 1 command execution")
		assert(t, executor.ExecutedCommands[0].Args[0] == "tmux", "Expected tmux command")
		assert(t, executor.ExecutedCommands[0].Args[1] == "switch", "Expected switch")
		assert(t, executor.ExecutedCommands[0].Args[3] == "mysession", "Expected session name")
	})

	t.Run("HasSession returns error if session name is empty", func(t *testing.T) {
		client := NewTmuxClient()
		executor := &MockCommandExecutor{}

		ok, err := client.HasSession(executor, "")
		assert(t, err != nil, "Expected error when session name is empty")
		assert(t, !ok, "Expected HasSession to return false")
	})

	t.Run("HasSession handles missing session gracefully", func(t *testing.T) {
		executor := &MockCommandExecutor{ReturnExitCode: 1, ReturnErr: errors.New("exit 1")}
		client := NewTmuxClient()

		ok, err := client.HasSession(executor, "nosession")

		assert(t, err == nil, "Expected no error for non-existent session")
		assert(t, !ok, "Expected HasSession to return false")
	})

	t.Run("HasSession returns true if session exists", func(t *testing.T) {
		executor := &MockCommandExecutor{ReturnExitCode: 0, ReturnErr: nil}
		client := NewTmuxClient()

		ok, err := client.HasSession(executor, "mysession")

		assert(t, err == nil, "Expected no error for existing session")
		assert(t, ok, "Expected HasSession to return true")
	})

	t.Run("NewSession returns error for missing required fields", func(t *testing.T) {
		client := NewTmuxClient()
		executor := &MockCommandExecutor{}

		err := client.NewSession(executor, "", "root", "win", "")
		assert(t, err != nil, "Expected error for empty session name")

		err = client.NewSession(executor, "sess", "", "win", "")
		assert(t, err != nil, "Expected error for empty session root")

		err = client.NewSession(executor, "sess", "root", "", "")
		assert(t, err != nil, "Expected error for empty window name")
	})

	t.Run("NewSession properly assembles the session", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := NewTmuxClient()

		err := client.NewSession(executor, "mysession", "/home/test", "main", "/project")
		assert(t, err == nil, "Expected no error")

		args := executor.ExecutedCommands[0].Args
		assert(t, args[0] == "tmux", "Expected tmux command")
		assert(t, args[1] == "new-session", "Expected new-session command")
		assert(t, args[2] == "-d", "Expected -d flag")
		assert(t, args[3] == "-s", "Expected -s flag")
		assert(t, args[4] == "mysession", "Expected session name")
		assert(t, args[5] == "-c", "Expected -c flag")
		assert(t, args[6] == "/home/test", "Expected session root")
		assert(t, args[7] == "-n", "Expected -n flag")
		assert(t, args[8] == "main", "Expected window name")
		assert(t, args[9] == "cd /project && exec $SHELL", "Expected window root")
	})

	t.Run("NewSession includes default root when windowRoot is empty", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := NewTmuxClient()

		err := client.NewSession(executor, "mysession", "/home/test", "main", "")
		assert(t, err == nil, "Expected no error")

		args := executor.ExecutedCommands[0].Args
		assert(t, args[0] == "tmux", "Expected tmux command")
		assert(t, args[1] == "new-session", "Expected new-session command")
		assert(t, args[2] == "-d", "Expected -d flag")
		assert(t, args[3] == "-s", "Expected -s flag")
		assert(t, args[4] == "mysession", "Expected session name")
		assert(t, args[5] == "-c", "Expected -c flag")
		assert(t, args[6] == "/home/test", "Expected session root")
		assert(t, args[7] == "-n", "Expected -n flag")
		assert(t, args[8] == "main", "Expected window name")
	})

	t.Run("SendKeys returns error if any input is empty", func(t *testing.T) {
		client := NewTmuxClient()
		executor := &MockCommandExecutor{}

		err := client.SendKeys(executor, "", "win", "ls")
		assert(t, err != nil, "Expected error for empty session name")

		err = client.SendKeys(executor, "sess", "", "ls")
		assert(t, err != nil, "Expected error for empty window name")

		err = client.SendKeys(executor, "sess", "win", "")
		assert(t, err != nil, "Expected error for empty keys")
	})

	t.Run("SendKeys properly assembles the keys", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := NewTmuxClient()

		err := client.SendKeys(executor, "mysession", "main", "ls")
		assert(t, err == nil, "Expected no error")

		args := executor.ExecutedCommands[0].Args
		assert(t, args[0] == "tmux", "Expected tmux command")
		assert(t, args[1] == "send-keys", "Expected send-keys command")
		assert(t, args[2] == "-t", "Expected -t flag")
		assert(t, args[3] == "mysession:main", "Expected session:window name")
		assert(t, args[4] == "ls", "Expected keys")
		assert(t, args[5] == "C-m", "Expected C-m")
	})

	t.Run("NewWindow returns error if required input is empty", func(t *testing.T) {
		client := NewTmuxClient()
		executor := &MockCommandExecutor{}

		err := client.NewWindow(executor, "", "/project", "window", "/root")
		assert(t, err != nil, "Expected error for empty session name")

		err = client.NewWindow(executor, "sess", "/project", "", "/root")
		assert(t, err != nil, "Expected error for empty window name")

		err = client.NewWindow(executor, "sess", "", "window", "/root")
		assert(t, err != nil, "Expected error for empty session root")
	})

	t.Run("NewWindow includes default root when windowRoot is empty", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := NewTmuxClient()

		err := client.NewWindow(executor, "sess", "/project", "main", "")
		assert(t, err == nil, "Expected no error")

		args := executor.ExecutedCommands[0].Args
		assert(t, args[0] == "tmux", "Expected tmux command")
		assert(t, args[1] == "new-window", "Expected new-window command")
		assert(t, args[len(args)-2] == "-c", "Expected -c flag before root")
		assert(t, args[len(args)-1] == "/project", "Expected session root used as fallback")
	})

	t.Run("NewWindow properly assembles the window", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := NewTmuxClient()

		err := client.NewWindow(executor, "mysession", "/home/test", "main", "/project")
		assert(t, err == nil, "Expected no error")

		args := executor.ExecutedCommands[0].Args
		assert(t, args[0] == "tmux", "Expected tmux command")
		assert(t, args[1] == "new-window", "Expected new-window command")
		assert(t, args[2] == "-d", "Expected -d flag")
		assert(t, args[3] == "-t", "Expected -t flag")
		assert(t, args[4] == "mysession", "Expected session name")
		assert(t, args[5] == "-n", "Expected -n flag")
		assert(t, args[6] == "main", "Expected window name")
		assert(t, args[7] == "-c", "Expected -c flag")
		assert(t, args[8] == "/project", "Expected window root")
	})
}

func Test_TmuxMultiplexer(t *testing.T) {
	t.Run("AttachProject returns error if project has no name", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		mu := NewTmuxMultiplexer(executor, nil)

		err := mu.AttachProject(&Project{Name: "", Template: &Template{Name: ""}})
		assert(t, err != nil, "Expected error when project has no name")
	})

	t.Run("AttachProject returns error if project template is nil", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		mu := NewTmuxMultiplexer(executor, nil)

		err := mu.AttachProject(&Project{Name: "foo", Template: nil})
		assert(t, err != nil, "Expected error when project has no template")
	})

	t.Run("AttachProject returns error if project template has no windows", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := &TmuxClient{}
		client.HasSession = func(e CommandExecutor, sessionName string) (bool, error) {
			return false, nil
		}

		mu := NewTmuxMultiplexer(executor, client)

		err := mu.AttachProject(&Project{Name: "foo", Template: &Template{Windows: []Window{}}})
		assert(t, err != nil, "Expected error when project template has no windows")
	})

	t.Run("AttachProject creates new session and attaches if it doesn't exist", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := &TmuxClient{}
		client.HasSession = func(e CommandExecutor, sessionName string) (bool, error) {
			return false, nil
		}
		client.NewSession = func(e CommandExecutor, sessionName, sessionRoot, windowName, windowRoot string) error {
			assert(t, sessionName == "foo", "Expected session name")
			assert(t, sessionRoot == "/home/test", "Expected session root")
			assert(t, windowName == "main", "Expected window name")
			assert(t, windowRoot == "/project", "Expected window root")
			return nil
		}

		client.IsInTmuxSession = func() bool {
			return false
		}

		client.AttachSession = func(e CommandExecutor, sessionName string) error {
			assert(t, sessionName == "foo", "Expected session name")
			return nil
		}

		mu := NewTmuxMultiplexer(executor, client)

		err := mu.AttachProject(&Project{Name: "foo", Template: &Template{Root: "/home/test", Windows: []Window{{Name: "main", Root: "/project"}}}})
		assert(t, err == nil, "Expected no error")
	})

	t.Run("AttachProject attaches to existing session if it exists", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := &TmuxClient{}
		client.HasSession = func(e CommandExecutor, sessionName string) (bool, error) {
			return true, nil
		}

		client.IsInTmuxSession = func() bool {
			return false
		}

		client.AttachSession = func(e CommandExecutor, sessionName string) error {
			assert(t, sessionName == "foo", "Expected session name")
			return nil
		}

		mu := NewTmuxMultiplexer(executor, client)

		err := mu.AttachProject(&Project{Name: "foo", Template: &Template{Root: "/home/test", Windows: []Window{{Name: "main", Root: "/project"}}}})
		assert(t, err == nil, "Expected no error")
	})

	t.Run("AttachProject creates new session and switches if it doesn't exist and is in tmux session", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := &TmuxClient{}
		client.HasSession = func(e CommandExecutor, sessionName string) (bool, error) {
			return false, nil
		}
		client.NewSession = func(e CommandExecutor, sessionName, sessionRoot, windowName, windowRoot string) error {
			assert(t, sessionName == "foo", "Expected session name")
			assert(t, sessionRoot == "/home/test", "Expected session root")
			assert(t, windowName == "main", "Expected window name")
			assert(t, windowRoot == "/project", "Expected window root")
			return nil
		}

		client.IsInTmuxSession = func() bool {
			return true
		}

		client.SwitchSession = func(e CommandExecutor, sessionName string) error {
			assert(t, sessionName == "foo", "Expected session name")
			return nil
		}

		mu := NewTmuxMultiplexer(executor, client)

		err := mu.AttachProject(&Project{Name: "foo", Template: &Template{Root: "/home/test", Windows: []Window{{Name: "main", Root: "/project"}}}})
		assert(t, err == nil, "Expected no error")
	})

	t.Run("AttachProject switches to existing session if it exists and is in tmux session", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := &TmuxClient{}
		client.HasSession = func(e CommandExecutor, sessionName string) (bool, error) {
			return true, nil
		}

		client.IsInTmuxSession = func() bool {
			return true
		}

		client.SwitchSession = func(e CommandExecutor, sessionName string) error {
			assert(t, sessionName == "foo", "Expected session name")
			return nil
		}

		mu := NewTmuxMultiplexer(executor, client)

		err := mu.AttachProject(&Project{Name: "foo", Template: &Template{Root: "/home/test", Windows: []Window{{Name: "main", Root: "/project"}}}})
		assert(t, err == nil, "Expected no error")
	})

	t.Run("AttachProject properly assembles a complex session", func(t *testing.T) {
		executor := &MockCommandExecutor{}
		client := &TmuxClient{}
		client.HasSession = func(e CommandExecutor, sessionName string) (bool, error) {
			return false, nil
		}
		client.NewSession = func(e CommandExecutor, sessionName, sessionRoot, windowName, windowRoot string) error {
			assert(t, sessionName == "foo", "Expected session name")
			assert(t, sessionRoot == "/home/test", "Expected session root")
			assert(t, windowName == "main", "Expected window name")
			assert(t, windowRoot == "/project", "Expected window root")
			return nil
		}

		client.NewWindow = func(e CommandExecutor, sessionName, sessionRoot, windowName, windowRoot string) error {
			assert(t, sessionName == "foo", "Expected session name")
			assert(t, sessionRoot == "/home/test", "Expected session root")
			assert(t, windowName == "baz", "Expected window name")
			assert(t, windowRoot == "", "Expected no window root")
			return nil
		}

		sendKeysCall := 0

		client.SendKeys = func(e CommandExecutor, sessionName, windowName, keys string) error {
			sendKeysCall++

			switch sendKeysCall {
			case 1:
				{
					assert(t, sessionName == "foo", "Expected session name")
					assert(t, windowName == "main", "Expected window name to be main is %s", windowName)
					assert(t, keys == "echo hello", "Expected keys")
				}
			case 2:
				{
					assert(t, sessionName == "foo", "Expected session name")
					assert(t, windowName == "baz", "Expected window name to be baz is %s", windowName)
					assert(t, keys == "echo hello", "Expected keys")
				}
			case 3:
				{
					assert(t, sessionName == "foo", "Expected session name")
					assert(t, windowName == "baz", "Expected window name to be baz is %s", windowName)
					assert(t, keys == "ls", "Expected keys")
				}
			}
			return nil
		}

		client.IsInTmuxSession = func() bool {
			return false
		}

		client.AttachSession = func(e CommandExecutor, sessionName string) error {
			assert(t, sessionName == "foo", "Expected session name")
			return nil
		}

		mu := NewTmuxMultiplexer(executor, client)

		err := mu.AttachProject(&Project{
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
		})

		assert(t, err == nil, "Expected no error")
	})
}
