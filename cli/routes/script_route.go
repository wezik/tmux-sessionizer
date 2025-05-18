package routes

import (
	"fmt"
	"os"
	"strings"
)

func ScriptRoute(args []string) {

	if len(args) == 0 {
		fmt.Println("Missing script command")
		os.Exit(1)
	}

	switch strings.ToLower(args[0]) {
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
}
