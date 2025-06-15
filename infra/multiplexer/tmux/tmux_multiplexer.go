package tmux

import (
	"fmt"
	"thop/dom/executor"
	"thop/dom/model/project"

	"github.com/dsnet/try"
)

type TmuxMultiplexer struct {
	client TmuxClient
}

type SessionName string

func New(executor executor.CommandExecutor) *TmuxMultiplexer {
	// creating tmux client here as it's just an inetrnal abstraction that should be hidden from clients
	client := NewTmuxClient(executor)
	return &TmuxMultiplexer{client: client}
}

// This function assumes that the project is validated before it's passed to it
func (m *TmuxMultiplexer) AttachProject(p *project.Project) error {
	sessionName := resolveSessionName(p)

	sessionExists, err := m.client.HasSession(sessionName)
	if err != nil {
		return err
	}

	if !sessionExists {
		sessionRoot := p.Template.Root
		mainWindow := p.Template.Windows[0]

		// first window gets created together with the session
		try.E(m.client.NewSession(sessionName, sessionRoot, mainWindow.Name, mainWindow.Root))

		for i, window := range p.Template.Windows {
			// main window is already created, so skip it
			if i != 0 {
				try.E(m.client.NewWindow(sessionName, sessionRoot, window.Name, window.Root))
			}

			for _, keys := range p.Template.Commands {
				try.E(m.client.SendKeys(sessionName, window.Name, keys))
			}

			for _, keys := range window.Commands {
				try.E(m.client.SendKeys(sessionName, window.Name, keys))
			}
		}

		fmt.Println("Session", sessionName, "created")
	}

	if m.client.IsInTmuxSession() {
		fmt.Println("Switching to", sessionName, "session")
		try.E(m.client.SwitchSession(sessionName))
	} else {
		fmt.Println("Attaching to", sessionName, "session")
		try.E(m.client.AttachSession(sessionName))
	}

	return nil
}

func resolveSessionName(p *project.Project) SessionName {
	if p.Template.Name == "" {
		return SessionName(p.Name)
	}

	return SessionName(p.Template.Name)
}
