package tmux

import (
	"errors"
	"fmt"
	. "thop/dom/model"
)

type TmuxMultiplexer struct {
	c TmuxClient
}

func NewTmuxMultiplexer(client TmuxClient) *TmuxMultiplexer {
	return &TmuxMultiplexer{c: client}
}

func (m *TmuxMultiplexer) AttachProject(p *Project) error {
	sessionName, err := resolveSessionName(p)
	if err != nil {
		return errors.New("failed to resolve session name")
	}

	sessionExists, err := m.c.HasSession(sessionName)
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
		err := m.c.NewSession(sessionName, sessionRoot, mainWindow.Name, mainWindow.Root)
		if err != nil {
			return err
		}

		for i, window := range p.Template.Windows {
			// main window is already created, so skip it
			if i != 0 {
				err := m.c.NewWindow(sessionName, sessionRoot, window.Name, window.Root)
				if err != nil {
					return err
				}
			}

			for _, keys := range p.Template.Commands {
				err := m.c.SendKeys(sessionName, window.Name, keys)
				if err != nil {
					return err
				}
			}

			for _, keys := range window.Commands {
				err := m.c.SendKeys(sessionName, window.Name, keys)
				if err != nil {
					return err
				}
			}
		}

		fmt.Println("Session", sessionName, "created")
	}

	if m.c.IsInTmuxSession() {
		fmt.Println("Switching to", sessionName, "session")
		err = m.c.SwitchSession(sessionName)
	} else {
		fmt.Println("Attaching to", sessionName, "session")
		err = m.c.AttachSession(sessionName)
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
