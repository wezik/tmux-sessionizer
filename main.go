package main

import (
	"thop/cmd"
	"thop/internal/executor"
	"thop/internal/fsystem"
	"thop/internal/selector"
	"thop/internal/service"
	"thop/internal/storage"
)

func main() {
	// userConfigDir := try.E1(os.UserConfigDir())
	//
	// configPath := filepath.Join(userConfigDir, "thop")
	// cfg := cfg.NewConfig(configPath)
	//
	// fs := fs.NewOsFileSystem()
	//
	// e := shell.NewCommandExecutor()
	// sl := fzf.NewFzfSelector(e)
	//
	// tc := tmux.NewTmuxClient(e)
	// mu := tmux.NewTmuxMultiplexer(tc)
	//
	// st := yaml.NewYamlStorage(cfg, fs)
	//
	// el := editor.NewShellEditorLauncher(cfg.GetEditor(), e)
	//
	// svc := service.NewService(sl, mu, st, el)
	executor := executor.ShellExecutor{}
	fsystem := fsystem.OsFileSystem{}

	svc := service.AppService{
		Selector:       selector.FzfSelector{E: executor},
		Multiplexer:    service.TmuxMultiplexer{E: executor},
		Storage:        storage.YamlStorage{},
		E:              executor,
	}

	cmd.AppService = &svc
	cmd.Execute()
}
