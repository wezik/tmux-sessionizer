package config

type Config struct {
	ConfigDir  string
	Editor     string
	InsideTmux bool
}

func (c *Config) GetConfigDir() string { return c.ConfigDir }
func (c *Config) GetEditor() string    { return c.Editor }
func (c *Config) IsInsideTmux() bool   { return c.InsideTmux }
