package cli

import (
	"os"
	"phopper/cli/routes"
	"phopper/domain/globals"
)

// TODO: think if it is a good idea to run on start, maybe should run on read err?
// or maybe it should act differently if storage is DB vs FS
func onStart() {
	globals.Get().Database.RunMigrations()
}

func Run() {
	onStart()
	args := os.Args[1:] // strip first arg (executable name)
	routes.MainRoute(args)
}
