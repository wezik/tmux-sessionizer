package session_template

const version = 1

type SessionTemplate struct {
	Version int `yaml:"version"`
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
