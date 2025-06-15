package selector

import (
	"thop/dom/problem"
)

const (
	ErrCancelled      problem.Key = "SELECTOR_CANCELLED"
	ErrSelectorFailed problem.Key = "SELECTOR_FAILED"
)

type Selector interface {
	SelectFrom(items []string, prompt string) (string, error)
}
