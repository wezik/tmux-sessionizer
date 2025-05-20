package routes

import (
	"fmt"
	"os"
	"phopper/domain/config"
	"phopper/domain/errors"
	"phopper/domain/globals"
	"phopper/domain/project/project_service"
	"strings"
)

const helpText = `Usage: phop [command]

Available commands:
  a, add, c, create     Create a new project in the current working directory
  d, delete, r, remove  Delete a project
  e, edit               Edit a project with the given editor, defaults to system default
  s, script             Manage scripts (see phop script help for more info)
  l, list               List projects
  h, help               Show this help
`

func MainRoute(args []string) {
	// no args defaults to project selection
	if len(args) == 0 {
		selectProject()
		return
	}

	switch strings.ToLower(args[0]) {
	case "a", "add", "c", "create":
		createProject()
	case "d", "delete", "r", "remove":
		deleteProject()
	case "e", "edit":
		editProject(args[1:])
	case "s", "script":
		ScriptRoute(args[1:])
	case "l", "list":
		selectProject()
	default:
		fmt.Print(helpText)
	}
}

func selectProject() {
	multiplexer := globals.Get().Multiplexer
	selector := globals.Get().Selector

	pjs := project_service.List()

	selected, err := selector.SelectFrom(pjs, "Select project to attach > ")
	if err != nil {
		return
	}

	multiplexer.AssembleAndAttach(selected)
}

func createProject() {
	cwd, err := os.Getwd()
	errors.EnsureNotNil(err, "Could not get current working directory")

	saved := project_service.Create(cwd)
	fmt.Println("Successfully created", saved.Session.Name, "template")
}

func deleteProject() {
	selector := globals.Get().Selector

	pjs := project_service.List()

	selected, err := selector.SelectFrom(pjs, "Select project to delete > ")
	if err != nil {
		return
	}

	project_service.Delete(selected.UUID)
	fmt.Println("Successfully deleted", selected.Session.Name, "template")
}

func editProject(args []string) {
	editor := config.GetDefaults().Editor
	if len(args) >= 1 {
		editor = args[0]
	}

	selector := globals.Get().Selector

	pjs := project_service.List()
	selected, err := selector.SelectFrom(pjs, "Select project to edit > ")
	if err != nil {
		return
	}

	project_service.Edit(selected, editor)
	fmt.Println("Successfully edited", selected.Session.Name, "template")
}
