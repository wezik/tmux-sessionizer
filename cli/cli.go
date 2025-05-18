package cli

import (
	"fmt"
	"os"
	"phopper/domain/project"
	"phopper/infra/repository"
	"phopper/infra/selector"
	"strings"
)

func Run() {
	args := os.Args[1:]

	repo := repository.FSProjectRepository{}
	selector := selector.FzfSelector{}

	if len(args) == 0 {
		cmd := project.ListAndSelectCommand{
			Repository: repo,
			Selector: selector,
		}
		project.ListAndSelect(cmd)
		os.Exit(0)
	}

	switch strings.ToLower(args[0]) {
	case "a", "add", "c", "create":
		path, err := os.Getwd()
		if err != nil {
			fmt.Println("Could not get current working directory")
			os.Exit(1)
		}
		cmd := project.CreateProjectCommand {
			Path: path,
			Repository: repo,
		}
		project.CreateProject(cmd)

	case "d", "delete", "r", "remove":
		cmd := project.ListAndDeleteCommand{
			Repository: repo,
			Selector: selector,
		}
		project.ListAndDelete(cmd)

	case "e", "edit":
		fmt.Println("TODO edit a project")

	case "s", "script":
		if (len(args) < 2) {
			fmt.Println("Missing script command")
			os.Exit(1)
		}

		fmt.Println("TODO script management")

		switch strings.ToLower(args[1]) {
		case "c", "create", "a", "add":
			fmt.Println("TODO create a script")

		case "d", "delete", "r", "remove":
			fmt.Println("TODO delete a script")

		case "l", "list":
			fmt.Println("TODO list scripts")
		default:
			fmt.Println("Unknown script command")
			os.Exit(1)
		}

	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
