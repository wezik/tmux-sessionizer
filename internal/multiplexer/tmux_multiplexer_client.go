package multiplexer

import (
	"fmt"
	"os"
	"os/exec"
	"thop/internal/executor"
	"thop/internal/problem"
	"thop/internal/types/command"
	"thop/internal/types/template"
	"thop/internal/types/window"
)

type TmuxClient interface {
	AttachSession(SessionName) error
	SwitchSession(SessionName) error
	HasSession(SessionName) (bool, error)
	NewSession(SessionName, template.Root, window.Name, window.Root) error
	NewWindow(SessionName, template.Root, window.Name, window.Root) error
	SendKeys(SessionName, window.Name, command.Command) error
}

type TmuxClientImpl struct {
	E executor.CommandExecutor
}

const (
	ErrFailedToAttachSession problem.Key = "TMUX_FAILED_TO_ATTACH_SESSION"
	ErrFailedToSwitchSession problem.Key = "TMUX_FAILED_TO_SWITCH_SESSION"
	ErrFailedToCheckSession  problem.Key = "TMUX_FAILED_TO_CHECK_SESSION"
	ErrFailedToCreateSession problem.Key = "TMUX_FAILED_TO_CREATE_SESSION"
	ErrFailedToCreateWindow  problem.Key = "TMUX_FAILED_TO_CREATE_WINDOW"
	ErrFailedToSendKeys      problem.Key = "TMUX_FAILED_TO_SEND_KEYS"
	ErrInvalidTemplateArgs   problem.Key = "TMUX_INVALID_TEMPLATE_ARGS"
)

func (c *TmuxClientImpl) AttachSession(session SessionName) error {
	if session == "" {
		return ErrInvalidTemplateArgs.WithMsg("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "attach", "-t", string(session))
	cmd.Stdin = os.Stdin // bind tmux session to terminal

	_, _, err := c.E.Execute(cmd)
	if err != nil {
		return ErrFailedToAttachSession.WithMsg(err.Error())
	}

	return nil
}

func (c *TmuxClientImpl) SwitchSession(session SessionName) error {
	if session == "" {
		return ErrInvalidTemplateArgs.WithMsg("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "switch", "-t", string(session))

	_, _, err := c.E.Execute(cmd)
	if err != nil {
		return ErrFailedToSwitchSession.WithMsg(err.Error())
	}

	return nil
}

func (c *TmuxClientImpl) HasSession(session SessionName) (bool, error) {
	if session == "" {
		return false, ErrInvalidTemplateArgs.WithMsg("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "has-session", "-t", string(session))

	_, exitCode, err := c.E.Execute(cmd)
	if err != nil {
		// exit code 1 means session does not exist
		if exitCode == 1 {
			return false, nil
		}
		return false, ErrFailedToCheckSession.WithMsg(err.Error())
	}

	return exitCode == 0, nil
}

func (c *TmuxClientImpl) NewSession(
	session SessionName,
	root template.Root,
	windowName window.Name,
	windowRoot window.Root,
) error {
	if anyEmpty(string(session), string(root), string(windowName)) {
		return ErrInvalidTemplateArgs.WithMsg("session, root and window name cannot be empty")
	}

	cmd := exec.Command("tmux", "new-session", "-d")
	cmd.Args = append(cmd.Args, "-s", string(session))
	cmd.Args = append(cmd.Args, "-c", string(root))
	cmd.Args = append(cmd.Args, "-n", string(windowName))

	if windowRoot != "" {
		// little hack to start first window at different root than session
		cmd.Args = append(cmd.Args, fmt.Sprintf("cd %s && exec $SHELL", windowRoot))
	}

	if _, _, err := c.E.Execute(cmd); err != nil {
		return ErrFailedToCreateSession.WithMsg(err.Error())
	}

	return nil
}

func (c *TmuxClientImpl) NewWindow(
	session SessionName,
	root template.Root,
	windowName window.Name,
	windowRoot window.Root,
) error {
	if anyEmpty(string(session), string(root), string(windowName)) {
		return ErrInvalidTemplateArgs.WithMsg("session, root and window name cannot be empty")
	}

	cmd := exec.Command("tmux", "new-window", "-d")
	cmd.Args = append(cmd.Args, "-t", string(session))
	cmd.Args = append(cmd.Args, "-n", string(windowName))

	if windowRoot != "" {
		cmd.Args = append(cmd.Args, "-c", string(windowRoot))
	} else {
		// in certain scenarios tmux will create window in working directory
		// instead of the session root, so specify it explicitly
		cmd.Args = append(cmd.Args, "-c", string(root))
	}

	if _, _, err := c.E.Execute(cmd); err != nil {
		return ErrFailedToCreateWindow.WithMsg(err.Error())
	}

	return nil

}

func (c *TmuxClientImpl) SendKeys(
	session SessionName,
	windowName window.Name,
	keys command.Command,
) error {
	if anyEmpty(string(session), string(windowName), string(keys)) {
		return ErrInvalidTemplateArgs.WithMsg("session, window name and keys cannot be empty")
	}

	cmd := exec.Command("tmux", "send-keys")

	// tmux needs combined name of session:window to send keys to
	cmd.Args = append(cmd.Args, "-t", fmt.Sprintf("%s:%s", session, windowName))
	cmd.Args = append(cmd.Args, string(keys))
	cmd.Args = append(cmd.Args, "C-m")

	if _, _, err := c.E.Execute(cmd); err != nil {
		return ErrFailedToSendKeys.WithMsg(err.Error())
	}

	return nil
}

func anyEmpty(s ...string) bool {
	for _, v := range s {
		if v == "" {
			return true
		}
	}
	return false
}
