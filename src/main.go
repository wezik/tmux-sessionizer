package main

import (
	"os"
	"phopper/src/app/cli"
	"phopper/src/domain/service"
)

func main() {
	svc := service.NewService(nil, nil)
	cli.NewCli(svc).Run(os.Args[1:])
}
