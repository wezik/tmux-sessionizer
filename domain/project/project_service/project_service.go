package project_service

import (
	"fmt"
	"phopper/domain/globals"
	"phopper/domain/project"
	"phopper/domain/project/session_template/template_service"
)

// Create

type CreateProjectCommand struct {
	Cwd string
}

func CreateProject(cmd CreateProjectCommand) project.Project {
	repo := globals.Get().ProjectRepository

	new_project := project.New(cmd.Cwd)
	saved := repo.SaveProject(new_project)

	fmt.Println("Successfully created", saved.Session.Name, "template")
	return saved
}

// List and select

func ListAndSelect() (project.Project, error) {
	selected, err := selectProject()
	// this means a search should just be canceled
	if err != nil { return project.Project{}, err }
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

	template_service.EditTemplate(editor, selected)
}

// Helper functions

func selectProject() (project.Project, error) {
	repo := globals.Get().ProjectRepository
	projects := repo.GetProjects()

	entries := make(map[string]project.Project)
	for _, project := range projects {
		entries[project.Session.Name] = project
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
