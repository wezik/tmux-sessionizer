package tmux

import (
	"fmt"
	"os"
	"thop/dom/executor"
	"thop/dom/model/command"
	"thop/dom/model/template"
	"thop/dom/model/window"
)

type TmuxClient interface {
	AttachSession(SessionName) error
	SwitchSession(SessionName) error
	HasSession(SessionName) (bool, error)
	NewSession(SessionName, template.Root, window.Name, window.Root) error
	NewWindow(SessionName, template.Root, window.Name, window.Root) error
	SendKeys(SessionName, window.Name, command.Command) error
	IsInTmuxSession() bool
}

type TmuxClientImpl struct {
	e executor.CommandExecutor
}

func NewTmuxClient(e executor.CommandExecutor) *TmuxClientImpl {
	return &TmuxClientImpl{e: e}
}

func (c *TmuxClientImpl) AttachSession(session SessionName) error {
	cmd := executor.Command("tmux", "attach", "-t", string(session))
	cmd.Stdin = os.Stdin // bind tmux session
	_, _, err := c.e.Execute(cmd)
	return err
}

func (c *TmuxClientImpl) SwitchSession(session SessionName) error {
	cmd := executor.Command("tmux", "switch", "-t", string(session))
	_, _, err := c.e.Execute(cmd)
	return err
}

func (c *TmuxClientImpl) HasSession(session SessionName) (bool, error) {
	cmd := executor.Command("tmux", "has-session", "-t", string(session))

	_, exitCode, err := c.e.Execute(cmd)
	if err != nil {
		// exit code 1 means session does not exist
		if exitCode == 1 {
			return false, nil
		}
		return false, err
	}

	return exitCode == 0, nil
}

func (c *TmuxClientImpl) NewSession(
	session SessionName,
	root template.Root,
	windowName window.Name,
	windowRoot window.Root,
) error {
	cmd := executor.Command("tmux", "new-session", "-d")
	cmd.Args = append(cmd.Args, "-s", string(session))
	cmd.Args = append(cmd.Args, "-c", string(root))
	cmd.Args = append(cmd.Args, "-n", string(windowName))

	if windowRoot != "" {
		// little hack to create first window at a different root than session
		cmd.Args = append(cmd.Args, fmt.Sprintf("cd %s && exec $SHELL", string(windowRoot)))
	}

	_, _, err := c.e.Execute(cmd)
	return err
}

func (c *TmuxClientImpl) NewWindow(
	session SessionName,
	root template.Root,
	windowName window.Name,
	windowRoot window.Root,
) (err error) {
	cmd := executor.Command("tmux", "new-window", "-d")
	cmd.Args = append(cmd.Args, "-t", string(session))
	cmd.Args = append(cmd.Args, "-n", string(windowName))

	if windowRoot != "" {
		cmd.Args = append(cmd.Args, "-c", string(windowRoot))
	} else {
		// in certain scenarios tmux will create window in working directory
		// instead of the session root, so specify it explicitly
		cmd.Args = append(cmd.Args, "-c", string(root))
	}

	_, _, err = c.e.Execute(cmd)
	return err
}

func (c *TmuxClientImpl) SendKeys(
	session SessionName,
	windowName window.Name,
	command command.Command,
) error {
	cmd := executor.Command("tmux", "send-keys")

	// tmux needs combined name of session:window to send keys to
	cmd.Args = append(cmd.Args, "-t", fmt.Sprintf("%s:%s", session, windowName))
	cmd.Args = append(cmd.Args, string(command))
	cmd.Args = append(cmd.Args, "C-m")

	_, _, err := c.e.Execute(cmd)
	return err
}

func (c *TmuxClientImpl) IsInTmuxSession() bool {
	// TMUX is set when tmux is running plain and simple
	return os.Getenv("TMUX") != ""
}
