package config

import (
	"os"
	. "phopper/src/domain/utils"
)

type Config struct {
	configDir string
	editor string
}

func NewConfig(configDir string) *Config {

	cfg := &Config{
		configDir: configDir,
	}

	err := os.MkdirAll(configDir, 0755)
	EnsureWithErr(err == nil, err)

	cfg.editor = defaultEditor

	// TODO: read editor from config file

	// in case editor is set it takes priority
	if editor := os.Getenv("EDITOR"); editor != "" {
		cfg.editor = editor
	}

	return cfg
}

func (c *Config) GetConfigDir() string { return c.configDir }
