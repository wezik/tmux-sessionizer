package cli

import (
	"fmt"
	"os"
	"phopper/domain/globals"
	"phopper/domain/project"
	"strings"
)

func setup() {
	globals.Get().Database.RunMigrations()
}

func Run() {
	setup()

	args := os.Args[1:]

	if len(args) == 0 {
		project.ListAndSelect()
		os.Exit(0)
	}

	switch strings.ToLower(args[0]) {
	case "a", "add", "c", "create":
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Could not get current working directory")
			os.Exit(1)
		}
		cmd := project.CreateProjectCommand {
			Cwd: cwd,
		}
		project.CreateProject(cmd)

	case "d", "delete", "r", "remove":
		project.ListAndDelete()

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
