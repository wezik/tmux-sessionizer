package routes

import (
	"fmt"
	"os"
	"phopper/domain/errors"
	"phopper/domain/globals"
	"phopper/domain/project/project_service"
	"strings"
)

const helpText = `Available commands:
  a, add, c, create     Create a new project in the current working directory
  d, delete, r, remove  Delete a project
  e, edit               Edit a project with the given editor, defaults to nano
  s, script             Manage scripts (see thop script help for more info)
  l, list               List projects
  h, help               Show this help

Usage: thop [command]
`

func MainRoute(args []string) {
	multiplexer := globals.Get().Multiplexer

	// no args defaults to project selection
	if len(args) == 0 {
		selected, err := project_service.ListAndSelect()
		if err != nil { return }
		multiplexer.AssembleAndAttach(selected)
		return
	}

	switch strings.ToLower(args[0]) {
	case "a", "add", "c", "create": {
		cwd, err := os.Getwd()
		errors.EnsureNotNil(err, "Could not get current working directory")

		project_service.CreateProject(project_service.CreateProjectCommand{Cwd: cwd})
	}
	
	case "d", "delete", "r", "remove":
		project_service.ListAndDelete()

	case "e", "edit": {
		var editor = "nano"
		if len(args) < 2 {
			fmt.Println("No editor specified, defaulting to nano")
		} else {
			editor = args[1]
		}
		project_service.ListAndEdit(editor)
	}

	case "s", "script":
		ScriptRoute(args[1:])

	case "l", "list":
		selected, err := project_service.ListAndSelect()
		if err != nil { return }
		multiplexer.AssembleAndAttach(selected)

	case "h", "help":
		fmt.Print(helpText)

	default:
		fmt.Println("Unknown command")
	}
}
