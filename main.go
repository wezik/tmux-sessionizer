package main

import (
	"os"
	"path/filepath"
	"thop/cmd"
	"thop/internal/config"
	"thop/internal/executor"
	"thop/internal/fsystem"
	"thop/internal/multiplexer"
	"thop/internal/selector"
	"thop/internal/service"
	"thop/internal/storage"
)

func main() {
	editor := os.Getenv("EDITOR")
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	configPath := filepath.Join(userConfigDir, "thop")
	tmuxSession := os.Getenv("TMUX")

	config := config.Config{
		ConfigDir: configPath,
		Editor:    editor,
	}

	executor := executor.ShellExecutor{}
	fsystem := fsystem.OsFileSystem{}

	svc := service.AppService{
		Selector: &selector.FzfProjectSelector{E: &executor},

		Multiplexer: &multiplexer.TmuxMultiplexer{
			ActiveTmuxSession: tmuxSession,
			Client:            &multiplexer.TmuxClientImpl{E: &executor},
		},

		Storage: &storage.YamlStorage{
			Config:     &config,
			FileSystem: &fsystem,
		},

		Config: &config,
		E:      &executor,
	}

	cmd.AppService = &svc
	cmd.Execute()
}
