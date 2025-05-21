package tmux

import (
	"phopper/domain/project"
)

type Tmux struct{}

var client = tmuxClient{}

func (t Tmux) AssembleAndAttach(project project.Project) {
	if !client.hasSession(project.Session.Name) {
		session := client.newSession(project)

		for i, window := range project.Session.Windows {
			// skip default window
			if i != 0 {
				client.newWindow(session, window)
			}

			for _, command := range project.Session.Commands {
				client.sendKeys(session, window.Name, command)
			}

			for _, command := range window.Commands {
				client.sendKeys(session, window.Name, command)
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

	if client.isInsideTmuxSession() {
		client.switchToSession(project.Session.Name)
	} else {
		client.attachToSession(project.Session.Name)
	}
}
