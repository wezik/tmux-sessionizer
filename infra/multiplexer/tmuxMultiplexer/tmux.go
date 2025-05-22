package tmuxMultiplexer

import (
	"phopper/domain/project"
)

type Tmux struct {
	Client *TmuxClient
}

func NewTmux() *Tmux {
	return &Tmux{
		Client: NewTmuxClient(),
	}
}

func (t *Tmux) AssembleAndAttach(project *project.Project) {
	if !t.Client.hasSession(project.Template.Name) {
		session := t.Client.newSession(project)

		for i, window := range project.Template.Windows {
			// skip default window
			if i != 0 {
				t.Client.newWindow(session, &window)
			}

			for _, command := range project.Template.Commands {
				t.Client.sendKeys(session, window.Name, command)
			}

			for _, command := range window.Commands {
				t.Client.sendKeys(session, window.Name, command)
			}
			//
			// for _, pane := range window.Panes {
			// 	client.newPane(project)
			//
			// 	for _, command := range pane.Commands {
			// 		client.sendKeys(project)
			// 	}
			// }
		}
	}

	if t.Client.isInsideTmuxSession() {
		t.Client.switchToSession(project.Template.Name)
	} else {
		t.Client.attachToSession(project.Template.Name)
	}
}
