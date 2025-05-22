package main

import (
	"os"
	"phopper/handlers"
)

func main() {
	handlers.NewProjectHandler().Run(os.Args[1:])
}
