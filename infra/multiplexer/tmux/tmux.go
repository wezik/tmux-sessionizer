package tmux

import (
	"os"
	"os/exec"
	"phopper/domain/errors"
	"phopper/domain/project"
)

type Tmux struct{}

func (t Tmux) AssembleAndAttach(project project.Project) {
	if !sessionExists(project.Session.Name) {
		createSession(project)
	}

	enterSession(project)
}

func sessionExists(session string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", session)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
}

func createSession(project project.Project) {
	cmd := exec.Command(
		"tmux",
		"new-session",
		"-d",
		"-s",
		project.Session.Name,
		"-c",
		project.Session.Path,
		"-n",
		"shell",
	)

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not create new session")
}

func enterSession(project project.Project) {
	// check if already in a tmux session
	sessionCmd := func() string {
		if len(os.Getenv("TMUX")) != 0 {
			return "switch"
		} else {
			return "attach"
		}
	}()

	cmd := exec.Command("tmux", sessionCmd, "-t", project.Session.Name)

	// bind to terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	errors.EnsureNotNil(err, "Could not attach to session")
}
