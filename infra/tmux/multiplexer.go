package tmux

import . "phopper/dom/model"
import . "phopper/dom/service"

type TmuxMultiplexer struct {
	e CommandExecutor
}

func NewTmuxMultiplexer(executor CommandExecutor) *TmuxMultiplexer {
	return &TmuxMultiplexer{e: executor}
}

func (m *TmuxMultiplexer) AttachProject(p *Project) error {
	panic("unimplemented")
}
