package tmux

import . "phopper/src/domain/model"
import . "phopper/src/domain/service"

type TmuxMultiplexer struct {
	e CommandExecutor
}

func NewTmuxMultiplexer(executor CommandExecutor) *TmuxMultiplexer {
	return &TmuxMultiplexer{e: executor}
}

func (m *TmuxMultiplexer) AttachProject(p *Project) error {
	panic("unimplemented")
}
