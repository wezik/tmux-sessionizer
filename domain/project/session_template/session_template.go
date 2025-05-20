package session_template

const version = "0.0.1"

type SessionTemplate struct {
	Version string `yaml:"version"` // used primarily by file system based storage
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
}

func New(path string) SessionTemplate {
	return SessionTemplate{
		Version: version,
		Name:    path, // name defaults to path (prob could be passed as an arg)
		Path:    path,
	}
}
