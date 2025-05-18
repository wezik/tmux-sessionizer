package selector

type Selector interface {
	ListAndSelect(entries []string, prompt string) string
}
