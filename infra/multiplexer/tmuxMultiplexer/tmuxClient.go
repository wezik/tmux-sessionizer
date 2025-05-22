package tmuxMultiplexer

import (
	"fmt"
	"os"
	"os/exec"
	"phopper/domain/errors"
	"phopper/domain/project"
	"phopper/domain/project/template"
	"phopper/domain/shell"
)

type TmuxClient struct{
	ShellRunner shell.Runner
}

func NewTmuxClient() *TmuxClient {
	return &TmuxClient{
		ShellRunner: shell.NewDefaultRunner(),
	}
}

func (t TmuxClient) newSession(p *project.Project) string {
	cmd := exec.Command("tmux", "new-session", "-d")
	cmd.Args = append(cmd.Args, "-s", p.Template.Name)
	cmd.Args = append(cmd.Args, "-c", p.Template.Root)
	cmd.Args = append(cmd.Args, "-n", p.Template.Windows[0].Name)

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not create new session")

	return p.Template.Name
}

func (_ TmuxClient) hasSession(session string) bool {
	cmd := exec.Command("tmux", "has-session")
	cmd.Args = append(cmd.Args, "-t", session)

	err := cmd.Run()
	if err != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
}

func (_ TmuxClient) switchToSession(session string) {
	cmd := exec.Command("tmux", "switch")
	cmd.Args = append(cmd.Args, "-t", session)

	// bind to terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not attach to a session")
}

func (_ TmuxClient) attachToSession(session string) {
	cmd := exec.Command("tmux", "attach")
	cmd.Args = append(cmd.Args, "-t", session)

	// bind to terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not attach to a session")
}

// returns true if current process is inside a tmux session
func (_ TmuxClient) isInsideTmuxSession() bool {
	return len(os.Getenv("TMUX")) != 0
}

func (_ TmuxClient) newWindow(session string, window *template.Window) {
	cmd := exec.Command("tmux", "new-window", "-d")
	cmd.Args = append(cmd.Args, "-t", session)
	cmd.Args = append(cmd.Args, "-n", window.Name)
	if window.Root != "" {
		cmd.Args = append(cmd.Args, "-c", window.Root)
	}

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not create new window")
}

func (_ TmuxClient) sendKeys(session string, window string, command string) {
	cmd := exec.Command("tmux", "send-keys")
	combinedId := fmt.Sprintf("%s:%s", session, window)
	cmd.Args = append(cmd.Args, "-t", combinedId)
	cmd.Args = append(cmd.Args, command)
	cmd.Args = append(cmd.Args, "C-m")

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not send keys")
}
