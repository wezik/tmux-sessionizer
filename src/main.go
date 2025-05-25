package main

import (
	"os"
	"path/filepath"
	"phopper/src/app/cli"
	"phopper/src/app/config"
	"phopper/src/domain/service"
	"phopper/src/infrastructure/fzf"
	"phopper/src/infrastructure/shell"
	"phopper/src/infrastructure/tmux"
	"phopper/src/infrastructure/yaml"
)

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig(filepath.Join(userConfigDir, "phopper"))

	e := shell.NewCommandExecutor()
	sl := fzf.NewFzfSelector(e)

	mu := tmux.NewTmuxMultiplexer(e)

	st := yaml.NewYamlStorage(cfg)

	svc := service.NewService(sl, mu, st)
	cli.NewCli(svc).Run(os.Args[1:])
}
