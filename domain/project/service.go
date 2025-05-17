package project

import (
	"fmt"
	"phopper/domain"
)

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

type ListAndSelectCommand struct {
	Repository domain.ProjectRepository
	Config domain.Config 
}

func ListAndSelect(cmd ListAndSelectCommand) {
	projects := cmd.Repository.GetAllProjects()

	if (cmd.Config.UseFzf()) {
		// TODO: use fzf
		fmt.Println("TODO: use fzf")
	} else {
		for _, project := range projects {
			fmt.Println(project)
		}
	}
}
