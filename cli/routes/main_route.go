package routes

import (
	"fmt"
	"os"
	"phopper/domain/project/project_service"
	"strings"
)

func MainRoute(args []string) {
	// no args defaults to project selection
	if len(args) == 0 {
		project_service.ListAndSelect()
		return
	}

	switch strings.ToLower(args[0]) {
	case "a", "add", "c", "create":
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Could not get current working directory")
			os.Exit(1)
		}

		cmd := project_service.CreateProjectCommand {
			Cwd: cwd,
		}

		project_service.CreateProject(cmd)
	
	case "d", "delete", "r", "remove":
		project_service.ListAndDelete()

	case "e", "edit":
		fmt.Println("TODO edit a project")

	case "s", "script":
		ScriptRoute(args[1:])

	case "l", "list":
		project_service.ListAndSelect()

	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
