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
	routes.MainRoute(os.Args[1:])
}
