package tmux

import "phopper/src/domain/model"

type TmuxMultiplexer struct{}

func NewTmuxMultiplexer() *TmuxMultiplexer {
	return &TmuxMultiplexer{}
}

func (m *TmuxMultiplexer) AttachProject(p *model.Project) error {
	panic("unimplemented")
}
