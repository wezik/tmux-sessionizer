package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"phopper/domain/errors"
	"phopper/domain/project"
	"phopper/domain/project/template"
)

type tmuxClient struct{}

func (_ tmuxClient) newSession(p project.Project) string {
	cmd := exec.Command("tmux", "new-session", "-d")
	cmd.Args = append(cmd.Args, "-s", p.Session.Name)
	cmd.Args = append(cmd.Args, "-c", p.Session.Root)
	cmd.Args = append(cmd.Args, "-n", p.Session.Windows[0].Name)

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not create new session")

	return p.Session.Name
}

func (_ tmuxClient) hasSession(session string) bool {
	cmd := exec.Command("tmux", "has-session")
	cmd.Args = append(cmd.Args, "-t", session)

	err := cmd.Run()
	if err != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
}

func (_ tmuxClient) switchToSession(session string) {
	cmd := exec.Command("tmux", "switch")
	cmd.Args = append(cmd.Args, "-t", session)

	// bind to terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not attach to a session")
}

func (_ tmuxClient) attachToSession(session string) {
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
func (_ tmuxClient) isInsideTmuxSession() bool {
	return len(os.Getenv("TMUX")) != 0
}

func (_ tmuxClient) newWindow(session string, window template.Window) {
	cmd := exec.Command("tmux", "new-window", "-d")
	cmd.Args = append(cmd.Args, "-t", session)
	cmd.Args = append(cmd.Args, "-n", window.Name)
	if window.Root != "" {
		cmd.Args = append(cmd.Args, "-c", window.Root)
	}

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not create new window")
}

func (_ tmuxClient) sendKeys(session string, window string, command string) {
	cmd := exec.Command("tmux", "send-keys")
	combinedId := fmt.Sprintf("%s:%s", session, window)
	cmd.Args = append(cmd.Args, "-t", combinedId)
	cmd.Args = append(cmd.Args, command)
	cmd.Args = append(cmd.Args, "C-m")

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not send keys")
}
