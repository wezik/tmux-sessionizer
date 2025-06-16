package multiplexer

import "thop/internal/types/project"

type Multiplexer interface {
	AttachProject(p project.Project) error
}
