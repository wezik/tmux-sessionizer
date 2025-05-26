package main

import (
	"os"
	"path/filepath"
	"phopper/cfg"
	"phopper/cli"
	"phopper/dom/service"
	"phopper/infra/fs"
	"phopper/infra/fzf"
	"phopper/infra/shell"
	"phopper/infra/tmux"
	"phopper/infra/yaml"
)

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	cfg := cfg.NewConfig(filepath.Join(userConfigDir, "phopper"))
	fs := fs.NewOsFileSystem()

	e := shell.NewCommandExecutor()
	sl := fzf.NewFzfSelector(e)

	mu := tmux.NewTmuxMultiplexer(e)

	st := yaml.NewYamlStorage(cfg, fs)

	svc := service.NewService(sl, mu, st)
	cli.NewCli(svc).Run(os.Args[1:])
}
