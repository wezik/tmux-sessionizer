package project

type Project struct {
	UUID string
	Name string
	Path string
}

func New(name string, path string) Project {
	return Project{
		Name: name,
		Path: path,
	}
}
