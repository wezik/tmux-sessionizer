package globals

import (
	"phopper/domain/database"
	"phopper/domain/database/repository"
	"phopper/domain/multiplexer"
	"phopper/domain/selector"
	// "phopper/infra/database/sqlite_database"
	"phopper/infra/database/yaml_database"
	"phopper/infra/multiplexer/tmux"
	"phopper/infra/selector/fzf_selector"
)

type Globals struct {
	// home-made DI
	Database          database.Database
	ProjectRepository repository.ProjectRepository
	Selector          selector.Selector
	Multiplexer       multiplexer.Multiplexer
}

func Get() Globals {
	// db := sqlite_database.SqliteDatabase{}
	db := yaml_database.YamlDatabase{}

	return Globals{
		Database:          db,
		ProjectRepository: db.GetProjectRepository(),
		Selector:          fzf_selector.FzfSelector{},
		Multiplexer:       tmux.Tmux{},
	}
}
