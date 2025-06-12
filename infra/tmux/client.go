package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	. "thop/dom/service"

	"github.com/dsnet/try"
)

type TmuxClient interface {
	AttachSession(sessionName string) error
	SwitchSession(sessionName string) error
	HasSession(sessionName string) (bool, error)
	NewSession(sessionName, sessionRoot, windowName, windowRoot string) error
	NewWindow(sessionName, sessionRoot, windowName, windowRoot string) error
	SendKeys(sessionName, windowName, keys string) error
	IsInTmuxSession() bool
}

type TmuxClientImpl struct {
	e CommandExecutor
}

func NewTmuxClient(e CommandExecutor) *TmuxClientImpl {
	return &TmuxClientImpl{e: e}
}

func (c *TmuxClientImpl) AttachSession(sessionName string) (err error) {
	defer try.HandleF(&err, func() {
		err = fmt.Errorf("tmux attach failed: %w", err)
	})

	if sessionName == "" {
		return errors.New("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "attach", "-t", sessionName)
	cmd.Stdin = os.Stdin // bind tmux session to terminal

	try.E2(c.e.Execute(cmd))

	return nil
}

func (c *TmuxClientImpl) SwitchSession(sessionName string) (err error) {
	defer try.HandleF(&err, func() {
		err = fmt.Errorf("tmux switch failed: %w", err)
	})

	if sessionName == "" {
		return errors.New("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "switch", "-t", sessionName)

	try.E2(c.e.Execute(cmd))

	return nil
}

func (c *TmuxClientImpl) HasSession(sessionName string) (bool, error) {
	if sessionName == "" {
		return false, errors.New("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "has-session", "-t", sessionName)

	_, exitCode, err := c.e.Execute(cmd)
	if err != nil {
		// exit code 1 means session does not exist
		if exitCode == 1 {
			return false, nil
		}
		return false, fmt.Errorf("tmux has-session failed: %w", err)
	}

	return exitCode == 0, nil
}

func (c *TmuxClientImpl) NewSession(sessionName, sessionRoot, windowName, windowRoot string) (err error) {
	defer try.HandleF(&err, func() {
		err = fmt.Errorf("tmux new-session failed: %w", err)
	})

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

	try.E2(c.e.Execute(cmd))

	return nil
}

func (c *TmuxClientImpl) NewWindow(sessionName, sessionRoot, windowName, windowRoot string) (err error) {
	defer try.HandleF(&err, func() {
		err = fmt.Errorf("tmux new-window failed: %w", err)
	})

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

	try.E2(c.e.Execute(cmd))

	return nil

}

func (c *TmuxClientImpl) SendKeys(sessionName, windowName, keys string) (err error) {
	defer try.HandleF(&err, func() {
		err = fmt.Errorf("tmux send-keys failed: %w", err)
	})

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

	try.E2(c.e.Execute(cmd))

	return nil
}

func (c *TmuxClientImpl) IsInTmuxSession() bool {
	// TMUX is set when tmux is running plain and simple
	return os.Getenv("TMUX") != ""
}
