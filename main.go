package main

import (
	"os"
	"path/filepath"
	"thop/cfg"
	"thop/cli"
	"thop/dom/service"
	"thop/infra/editor"
	"thop/infra/fs"
	"thop/infra/fzf"
	"thop/infra/shell"
	"thop/infra/tmux"
	"thop/infra/yaml"
)

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	cfg := cfg.NewConfig(filepath.Join(userConfigDir, "thop"))
	fs := fs.NewOsFileSystem()

	e := shell.NewCommandExecutor()
	sl := fzf.NewFzfSelector(e)

	tc := tmux.NewTmuxClient(e)
	mu := tmux.NewTmuxMultiplexer(tc)

	st := yaml.NewYamlStorage(cfg, fs)

	el := editor.NewShellEditorLauncher(cfg.GetEditor(), e)

	svc := service.NewService(sl, mu, st, el)
	cli.NewCli(svc).Run(os.Args[1:])
}
