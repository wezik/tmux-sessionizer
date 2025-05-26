package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	. "phopper/dom/model"
	. "phopper/dom/service"
)

type TmuxMultiplexer struct {
	e CommandExecutor
}

func NewTmuxMultiplexer(executor CommandExecutor) *TmuxMultiplexer {
	return &TmuxMultiplexer{e: executor}
}

func (m *TmuxMultiplexer) AttachProject(p *Project) error {
	sessionName, err := resolveSessionName(p)
	if err != nil {
		return err
	}

	sessionExists, err := m.sessionExists(p)
	if err != nil {
		return err
	}

	if !sessionExists {
		err := m.newSession(p)
		if err != nil {
			return err
		}

		for _, window := range p.Template.Windows[1:] {
			err := m.newWindow(sessionName, &window, p.Template.Root)
			if err != nil {
				return err
			}

			for _, keys := range p.Template.Commands {
				err := m.sendKeys(sessionName, window.Name, keys)
				if err != nil {
					return err
				}
			}

			for _, keys := range window.Commands {
				err := m.sendKeys(sessionName, window.Name, keys)
				if err != nil {
					return err
				}
			}
		}

		fmt.Println("Session", sessionName, "created")
	}

	if isInsideTmuxSession() {
		fmt.Println("Switching to", sessionName, "session")
		err = m.switchToSession(p)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Attaching to", sessionName, "session")
		err = m.attachToSession(p)
		if err != nil {
			return err
		}
	}

	return err
}

func (m *TmuxMultiplexer) newSession(p *Project) error {
	sessionName, err := resolveSessionName(p)
	if err != nil {
		return err
	}

	templateRoot := p.Template.Root
	if templateRoot == "" {
		return errors.New("template root cannot be empty")
	}

	if len(p.Template.Windows) == 0 {
		return errors.New("template windows cannot be empty")
	}

	mainWindow := p.Template.Windows[0]

	windowName := mainWindow.Name
	if windowName == "" {
		return errors.New("window name cannot be empty")
	}

	cmd := exec.Command("tmux", "new-session", "-d")
	cmd.Args = append(cmd.Args, "-s", sessionName)
	cmd.Args = append(cmd.Args, "-c", templateRoot)
	cmd.Args = append(cmd.Args, "-n", windowName)

	_, _, err = m.e.Execute(cmd)
	if err != nil {
		return err
	}

	for _, keys := range p.Template.Commands {
		err := m.sendKeys(sessionName, windowName, keys)
		if err != nil {
			return err
		}
	}

	for _, keys := range mainWindow.Commands {
		err := m.sendKeys(sessionName, windowName, keys)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *TmuxMultiplexer) newWindow(session string, w *Window, templateRoot string) error {
	cmd := exec.Command("tmux", "new-window", "-d")
	cmd.Args = append(cmd.Args, "-t", session)

	if w.Name == "" {
		return errors.New("window name cannot be empty")
	}
	cmd.Args = append(cmd.Args, "-n", w.Name)

	if w.Root != "" {
		cmd.Args = append(cmd.Args, "-c", w.Root)
	} else {
		cmd.Args = append(cmd.Args, "-c", templateRoot)
	}

	_, _, err := m.e.Execute(cmd)
	return err
}

func (m *TmuxMultiplexer) sendKeys(session, window, keys string) error {
	cmd := exec.Command("tmux", "send-keys")
	combinedName := fmt.Sprintf("%s:%s", session, window)
	cmd.Args = append(cmd.Args, "-t", combinedName)
	cmd.Args = append(cmd.Args, keys)
	cmd.Args = append(cmd.Args, "C-m")

	_, _, err := m.e.Execute(cmd)
	return err
}

func (m *TmuxMultiplexer) attachToSession(p *Project) error {
	sessionName, err := resolveSessionName(p)
	if err != nil {
		return err
	}

	cmd := exec.Command("tmux", "attach", "-t", sessionName)
	cmd.Stdin = os.Stdin
	_, _, err = m.e.Execute(cmd)

	return err
}

func (m *TmuxMultiplexer) switchToSession(p *Project) error {
	sessionName, err := resolveSessionName(p)
	if err != nil {
		return err
	}

	cmd := exec.Command("tmux", "switch", "-t", sessionName)

	_, _, err = m.e.Execute(cmd)
	return err
}

func (m *TmuxMultiplexer) sessionExists(p *Project) (bool, error) {
	cmd := exec.Command("tmux", "has-session")

	sessionName, err := resolveSessionName(p)
	if err != nil {
		return false, err
	}

	cmd.Args = append(cmd.Args, "-t", sessionName)

	_, exitCode, err := m.e.Execute(cmd)
	if err != nil {
		if exitCode == 1 {
			return false, nil
		}
		return false, err
	}

	return exitCode == 0, err
}

func isInsideTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

func resolveSessionName(p *Project) (string, error) {
	if p.Template.Name == "" {
		if p.Name == "" {
			return "", errors.New("project name cannot be empty")
		}
		return p.Name, nil
	}

	return p.Template.Name, nil
}
