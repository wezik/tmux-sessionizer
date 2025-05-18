package globals

import (
	"phopper/domain"
	"phopper/infra/selector/fzf_selector"
	"phopper/infra/storage/sqlite_database"
)

type Globals struct {
	Database domain.Database
	ProjectRepository domain.ProjectRepository
	Selector domain.Selector
}

func Get() Globals {
	database := sqlite_database.SqliteDatabase{}

	return Globals{
		Database: database,
		ProjectRepository: database.GetProjectRepository(),
		Selector: fzf_selector.FzfSelector{},
	}
}

