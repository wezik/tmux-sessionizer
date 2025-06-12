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

	"github.com/dsnet/try"
)

func main() {
	userConfigDir := try.E1(os.UserConfigDir())

	configPath := filepath.Join(userConfigDir, "thop")
	cfg := cfg.NewConfig(configPath)

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
