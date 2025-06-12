package tmux

import (
	"errors"
	"fmt"
	. "thop/dom/model"

	"github.com/dsnet/try"
)

type TmuxMultiplexer struct {
	c TmuxClient
}

func NewTmuxMultiplexer(client TmuxClient) *TmuxMultiplexer {
	return &TmuxMultiplexer{c: client}
}

func (m *TmuxMultiplexer) AttachProject(p *Project) (err error) {
	defer try.Handle(&err)

	sessionName := try.E1(resolveSessionName(p))

	sessionExists := try.E1(m.c.HasSession(sessionName))

	if !sessionExists {
		if len(p.Template.Windows) == 0 {
			return errors.New("project template needs at least one window to be created")
		}

		sessionRoot := p.Template.Root
		mainWindow := p.Template.Windows[0]

		// first window gets created together with the session
		try.E(m.c.NewSession(sessionName, sessionRoot, mainWindow.Name, mainWindow.Root))

		for i, window := range p.Template.Windows {
			// main window is already created, so skip it
			if i != 0 {
				try.E(m.c.NewWindow(sessionName, sessionRoot, window.Name, window.Root))
			}

			for _, keys := range p.Template.Commands {
				try.E(m.c.SendKeys(sessionName, window.Name, keys))
			}

			for _, keys := range window.Commands {
				try.E(m.c.SendKeys(sessionName, window.Name, keys))
			}
		}

		fmt.Println("Session", sessionName, "created")
	}

	if m.c.IsInTmuxSession() {
		fmt.Println("Switching to", sessionName, "session")
		try.E(m.c.SwitchSession(sessionName))
	} else {
		fmt.Println("Attaching to", sessionName, "session")
		try.E(m.c.AttachSession(sessionName))
	}

	return nil
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
