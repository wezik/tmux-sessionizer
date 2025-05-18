package project

import (
	"fmt"
	"phopper/domain"
	"phopper/domain/globals"
)

// Create

type CreateProjectCommand struct {
	Cwd string
}

func CreateProject(cmd CreateProjectCommand) domain.TmuxProject {
	project := domain.TmuxProject{
		Name: cmd.Cwd,
		Path: cmd.Cwd,
	}

	repo := globals.Get().ProjectRepository
	saved := repo.SaveProject(project)

	fmt.Println("Created project:", saved)
	return saved
}

// List and select

func ListAndSelect() {
	selected := selectProject()
	fmt.Println("Selected:", selected)
}

// List and delete

func ListAndDelete() {
	selected := selectProject()
	repo := globals.Get().ProjectRepository
	repo.DeleteProject(selected.UUID)
}

// Helper functions

func selectProject() domain.TmuxProject {
	repo := globals.Get().ProjectRepository
	projects := repo.GetProjects()

	entries := make(map[string]domain.TmuxProject)
	for _, project := range projects {
		entries[project.Name] = project
	}

	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}

	selector := globals.Get().Selector
	selectedStr := selector.ListAndSelect(keys, "Select project > ")

	return entries[selectedStr]
}
