package multiplexer

import (
	"thop/dom/model/project"
	"thop/dom/problem"
)

const (
	ErrFailedToAttach = problem.Key("MULTIPLEXER_FAILED_TO_ATTACH")
)

type Multiplexer interface {
	AttachProject(p *project.Project) error
}
