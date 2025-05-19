package multiplexer

import "phopper/domain/project"

type Multiplexer interface {
	AssembleAndAttach(project project.Project)
}
