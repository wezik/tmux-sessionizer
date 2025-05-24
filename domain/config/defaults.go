package config

type Config struct {
	Editor string
}

// TODO this should be loaded from some sort of config file
const settingsEditor = "nvim"

func GetDefaults() Config {
	editor := settingsEditor
	if editor == "" {
		editor = defaultEditor
	}

	return Config{
		Editor: editor,
	}
}
