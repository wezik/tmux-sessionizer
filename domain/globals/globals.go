package globals

import (
	"phopper/domain/storage"
	"phopper/domain/storage/repository"
	"phopper/domain/multiplexer"
	"phopper/domain/selector"
	"phopper/infra/storage/yaml_storage"
	"phopper/infra/multiplexer/tmux"
	"phopper/infra/selector/fzf_selector"
)

type Globals struct {
	// DI at home
	Database          storage.Storage
	ProjectRepository repository.ProjectRepository
	Selector          selector.Selector
	Multiplexer       multiplexer.Multiplexer
}

func Get() Globals {
	storage := yaml_storage.YamlStorage{}

	return Globals{
		Database:          storage,
		ProjectRepository: storage.GetProjectRepository(),
		Selector:          fzf_selector.FzfSelector{},
		Multiplexer:       tmux.Tmux{},
	}
}
