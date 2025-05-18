package cli

import (
	"os"
	"phopper/cli/routes"
	"phopper/domain/globals"
)

func onStart() {
	globals.Get().Database.RunMigrations()
}

func Run() {
	onStart()
	args := os.Args[1:] // strip first arg (executable name)
	routes.MainRoute(args)
}
