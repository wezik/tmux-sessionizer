package cfg

type Editor string
type ConfigPath string
type Config struct {
	ConfigPath ConfigPath
	Editor     Editor
}

func New(configPath ConfigPath, editor Editor) *Config {
	return &Config{
		ConfigPath: configPath,
		Editor:     editor,
	}
}

func (c ConfigPath) String() string {
	return string(c)
}

func (c Editor) String() string {
	return string(c)
}
