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

func ListAndSelect() (project.Project, error) {
	selected, err := selectProject()
	// this means a search should just be canceled
	if err != nil { return project.Project{}, err }

	fmt.Println("Selected:", selected)
	return selected, nil
}

// List and delete

func ListAndDelete() {
	selected, err := selectProject()
	// this means a search should just be canceled
	if err != nil { return }

	repo := globals.Get().ProjectRepository
	repo.DeleteProject(selected.UUID)
}

// List and edit

func ListAndEdit(editor string) {
	selected, err := selectProject()
	// this means a search should just be canceled
	if err != nil { return }
	fmt.Println("Editing:", selected)
	fmt.Println("TODO: create temp file and open it in editor, after save and pull changes")
}

// Helper functions

func selectProject() (project.Project, error) {
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

	selectedKey, err := selector.ListAndSelect(keys, "Select project > ")
	if err != nil { return project.Project{}, err }

	return entries[selectedKey], nil
}
