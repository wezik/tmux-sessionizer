package selector

import "phopper/domain/project"

type Selector interface {
	SelectProject(entries []project.Project, prompt string) (*project.Project, error)
}
