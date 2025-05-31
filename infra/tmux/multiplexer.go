package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	. "thop/dom/model"
	. "thop/dom/service"
)

type TmuxMultiplexer struct {
	e CommandExecutor
	c *TmuxClient
}

type TmuxClient struct {
	AttachSession   func(e CommandExecutor, sessionName string) error
	SwitchSession   func(e CommandExecutor, sessionName string) error
	HasSession      func(e CommandExecutor, sessionName string) (bool, error)
	NewSession      func(e CommandExecutor, sessionName, sessionRoot, windowName, windowRoot string) error
	NewWindow       func(e CommandExecutor, sessionName, sessionRoot, windowName, windowRoot string) error
	SendKeys        func(e CommandExecutor, sessionName, windowName, keys string) error
	IsInTmuxSession func() bool
}

func NewTmuxMultiplexer(commandExecutor CommandExecutor, client *TmuxClient) *TmuxMultiplexer {
	return &TmuxMultiplexer{e: commandExecutor, c: client}
}

func NewTmuxClient() *TmuxClient {
	return &TmuxClient{
		AttachSession:   attachSession,
		SwitchSession:   switchSession,
		HasSession:      hasSession,
		NewSession:      newSession,
		NewWindow:       newWindow,
		SendKeys:        sendKeys,
		IsInTmuxSession: isInTmuxSession,
	}
}

func attachSession(e CommandExecutor, sessionName string) error {
	if sessionName == "" {
		return errors.New("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "attach", "-t", sessionName)
	cmd.Stdin = os.Stdin // bind tmux session to terminal
	if _, _, err := e.Execute(cmd); err != nil {
		return errors.New("tmux attach failed")
	}

	return nil
}

func switchSession(e CommandExecutor, sessionName string) error {
	if sessionName == "" {
		return errors.New("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "switch", "-t", sessionName)

	if _, _, err := e.Execute(cmd); err != nil {
		return errors.New("tmux switch failed")
	}

	return nil
}

func hasSession(e CommandExecutor, sessionName string) (bool, error) {
	if sessionName == "" {
		return false, errors.New("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "has-session", "-t", sessionName)

	_, exitCode, err := e.Execute(cmd)
	if err != nil {
		// exit code 1 means session does not exist
		if exitCode == 1 {
			return false, nil
		}
		return false, errors.New("tmux has-session failed")
	}

	return exitCode == 0, nil
}

func newSession(e CommandExecutor, sessionName, sessionRoot, windowName, windowRoot string) error {
	if sessionName == "" {
		return errors.New("session name cannot be empty")
	}

	if sessionRoot == "" {
		return errors.New("session root cannot be empty")
	}

	if windowName == "" {
		return errors.New("window name cannot be empty")
	}

	cmd := exec.Command("tmux", "new-session", "-d")
	cmd.Args = append(cmd.Args, "-s", sessionName)
	cmd.Args = append(cmd.Args, "-c", sessionRoot)
	cmd.Args = append(cmd.Args, "-n", windowName)

	if windowRoot != "" {
		// little hack to start session at the root but modify the directory for first window
		cmd.Args = append(cmd.Args, fmt.Sprintf("cd %s && exec $SHELL", windowRoot))
	}

	if _, _, err := e.Execute(cmd); err != nil {
		return errors.New("tmux new-session failed")
	}

	return nil
}

func newWindow(e CommandExecutor, sessionName, sessionRoot, windowName, windowRoot string) error {
	if sessionName == "" {
		return errors.New("session name cannot be empty")
	}

	if windowName == "" {
		return errors.New("window name cannot be empty")
	}

	if sessionRoot == "" {
		return errors.New("session root cannot be empty")
	}

	cmd := exec.Command("tmux", "new-window", "-d")
	cmd.Args = append(cmd.Args, "-t", sessionName)
	cmd.Args = append(cmd.Args, "-n", windowName)

	if windowRoot != "" {
		cmd.Args = append(cmd.Args, "-c", windowRoot)
	} else {
		// in certain scenarios tmux will create window in working directory
		// instead of the session root, so specify it explicitly
		cmd.Args = append(cmd.Args, "-c", sessionRoot)
	}

	if _, _, err := e.Execute(cmd); err != nil {
		return errors.New("tmux new-window failed")
	}

	return nil

}

func sendKeys(e CommandExecutor, sessionName, windowName, keys string) error {
	if sessionName == "" {
		return errors.New("session name cannot be empty")
	}

	if windowName == "" {
		return errors.New("window name cannot be empty")
	}

	if keys == "" {
		// why are you sending empty keys?
		return errors.New("keys cannot be empty")
	}

	cmd := exec.Command("tmux", "send-keys")

	// tmux needs combined name of session:window to send keys to
	cmd.Args = append(cmd.Args, "-t", fmt.Sprintf("%s:%s", sessionName, windowName))
	cmd.Args = append(cmd.Args, keys)
	cmd.Args = append(cmd.Args, "C-m")

	if _, _, err := e.Execute(cmd); err != nil {
		return errors.New("tmux send-keys failed")
	}

	return nil

}

func isInTmuxSession() bool {
	// TMUX is set when tmux is running plain and simple
	return os.Getenv("TMUX") != ""
}

func (m *TmuxMultiplexer) AttachProject(p *Project) error {
	sessionName, err := resolveSessionName(p)
	if err != nil {
		return errors.New("failed to resolve session name")
	}

	sessionExists, err := m.c.HasSession(m.e, sessionName)
	if err != nil {
		return err
	}

	if !sessionExists {
		if len(p.Template.Windows) == 0 {
			return errors.New("project template needs at least one window to be created")
		}

		sessionRoot := p.Template.Root
		mainWindow := p.Template.Windows[0]

		// first window gets created together with the session
		err := m.c.NewSession(m.e, sessionName, sessionRoot, mainWindow.Name, mainWindow.Root)
		if err != nil {
			return err
		}

		for i, window := range p.Template.Windows {
			// main window is already created, so skip it
			if i != 0 {
				err := m.c.NewWindow(m.e, sessionName, sessionRoot, window.Name, window.Root)
				if err != nil {
					return err
				}
			}

			for _, keys := range p.Template.Commands {
				err := m.c.SendKeys(m.e, sessionName, window.Name, keys)
				if err != nil {
					return err
				}
			}

			for _, keys := range window.Commands {
				err := m.c.SendKeys(m.e, sessionName, window.Name, keys)
				if err != nil {
					return err
				}
			}
		}

		fmt.Println("Session", sessionName, "created")
	}

	if m.c.IsInTmuxSession() {
		fmt.Println("Switching to", sessionName, "session")
		err = m.c.SwitchSession(m.e, sessionName)
	} else {
		fmt.Println("Attaching to", sessionName, "session")
		err = m.c.AttachSession(m.e, sessionName)
	}

	return err
}

func resolveSessionName(p *Project) (string, error) {
	if p.Template == nil {
		return "", errors.New("project template cannot be nil")
	}

	if p.Template.Name == "" {
		if p.Name == "" {
			return "", errors.New("project name cannot be empty")
		}
		return p.Name, nil
	}

	return p.Template.Name, nil
}
