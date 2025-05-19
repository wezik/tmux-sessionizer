package project

type Project struct {
	UUID string
	Session SessionTemplate
}

type SessionTemplate struct {
	// keeping the name to a session as it's used primarily to differentiate between sessions
	// maybe it's a good idea to default to a project name, and make it configurable in the future
	Name string
	Path string
}

func New(name string, path string) Project {
	return Project{
		Session: SessionTemplate{Name: name, Path: path},
	}
}
