package project

import (
	"fmt"
	"phopper/domain"
)

// Create

type CreateProjectCommand struct {
	Repository domain.ProjectRepository
	Path string
}

func CreateProject(cmd CreateProjectCommand) domain.TmuxProject {
	project := domain.TmuxProject{
		Name: cmd.Path,
		Path: cmd.Path,
	}
	saved := cmd.Repository.SaveProject(project)

	fmt.Println("Created project:", saved)
	return saved
}

// List and select

type ListAndSelectCommand struct {
	Repository domain.ProjectRepository
	Selector domain.Selector
}

func ListAndSelect(cmd ListAndSelectCommand) {
	selected := selectProject(cmd.Repository, cmd.Selector)
	fmt.Println("Selected:", selected)
}

// List and delete

type ListAndDeleteCommand struct {
	Repository domain.ProjectRepository
	Selector domain.Selector
}

func ListAndDelete(cmd ListAndDeleteCommand) {
	selected := selectProject(cmd.Repository, cmd.Selector)
	cmd.Repository.DeleteProject(selected.UUID)
}

// Helper functions

func selectProject(repository domain.ProjectRepository, selector domain.Selector) domain.TmuxProject {
	projects := repository.GetAllProjects()

	stringifiedProjects := make([]string, len(projects))
	for i, project := range projects {
		stringifiedProjects[i] = project.String()
	}

	selectedStr := selector.ListAndSelect(stringifiedProjects, "Select project")
	return domain.TmuxProjectFromString(selectedStr)
}
