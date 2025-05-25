package main

import (
	"os"
	"phopper/src/app/cli"
	"phopper/src/domain/service"
	"phopper/src/infrastructure/fzf"
	"phopper/src/infrastructure/tmux"
	"phopper/src/infrastructure/yaml"
)

func main() {
	sl := fzf.NewFzfSelector()
	mu := tmux.NewTmuxMultiplexer()
	st := yaml.NewYamlStorage()
	svc := service.NewService(sl, mu, st)
	cli.NewCli(svc).Run(os.Args[1:])
}
