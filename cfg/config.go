package cfg

import (
	"os"
	. "phopper/dom/utils"
)

type Config interface {
	GetConfigDir() string
	GetEditor() string
}

type ConfigImpl struct {
	configDir string
	editor    string
}

func NewConfig(configDir string) *ConfigImpl {

	cfg := &ConfigImpl{
		configDir: configDir,
	}

	err := os.MkdirAll(configDir, 0755)
	EnsureWithErr(err == nil, err)

	cfg.editor = defaultEditor

	// in case editor is set it takes priority
	if editor := os.Getenv("EDITOR"); editor != "" {
		cfg.editor = editor
	}

	// TODO: override editor with config file?

	return cfg
}

func (c *ConfigImpl) GetConfigDir() string { return c.configDir }

func (c *ConfigImpl) GetEditor() string { return c.editor }
