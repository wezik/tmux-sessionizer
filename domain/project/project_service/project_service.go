package project_service

import (
	"fmt"
	"phopper/domain/globals"
	"phopper/domain/project"
)

// Create

type CreateProjectCommand struct {
	Cwd string
}

func CreateProject(cmd CreateProjectCommand) project.Project {
	repo := globals.Get().ProjectRepository

	new_project := project.New(cmd.Cwd, cmd.Cwd)
	saved := repo.SaveProject(new_project)

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

func selectProject() project.Project {
	repo := globals.Get().ProjectRepository
	projects := repo.GetProjects()

	entries := make(map[string]project.Project)
	for _, project := range projects {
		entries[project.Name] = project
	}

	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}

	selector := globals.Get().Selector
	selectedKey := selector.ListAndSelect(keys, "Select project > ")

	return entries[selectedKey]
}
