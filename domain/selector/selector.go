package selector

import "phopper/domain/project"

type Selector interface {
	// if this function returns an error, treat it as a canceled selection
	SelectFrom(entries []project.Project, prompt string) (project.Project, error)
}
