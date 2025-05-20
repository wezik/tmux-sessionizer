package cli

import (
	"os"
	"phopper/cli/routes"
)

func Run() {
	args := os.Args[1:] // strip first arg (executable name)
	routes.MainRoute(args)
}
