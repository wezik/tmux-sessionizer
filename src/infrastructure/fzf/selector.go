package fzf

type FzfSelector struct{}

func NewFzfSelector() *FzfSelector {
	return &FzfSelector{}
}

func (s *FzfSelector) SelectFrom(items []string) (string, error) {
	panic("unimplemented")
}
