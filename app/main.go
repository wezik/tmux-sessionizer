package main

import (
	"os"
	"path/filepath"
	"thop/app/cli"
	"thop/dom/cfg"
	"thop/dom/problem"
	"thop/dom/service"
	"thop/infra/executor/shell"
	"thop/infra/fsystem"
	"thop/infra/multiplexer/tmux"
	"thop/infra/selector/fzf"
	"thop/infra/storage/yaml"
)

const (
	ErrFailedToLoadConfigDir problem.Key = "THOP_FAILED_TO_LOAD_CONFIG_DIR"
)

func main() {
	cfg, err := loadConfig()
	handleErr(err)

	fs := fsystem.NewOsFileSystem()

	executor := shell.New()

	selector := fzf.New(executor)
	multiplexer := tmux.New(executor)

	storage := yaml.New(cfg, fs)

	svc := service.New(selector, multiplexer, storage, executor)

	cli := cli.New(svc)

	handleErr(cli.Run(os.Args[1:]))
}

func loadConfig() (*cfg.Config, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, ErrFailedToLoadConfigDir.WithMessage(err.Error())
	}

	configPath := cfg.ConfigPath(filepath.Join(userConfigDir, "thop"))
	editor := cfg.Editor(os.Getenv("EDITOR"))

	return cfg.New(configPath, editor), nil
}
