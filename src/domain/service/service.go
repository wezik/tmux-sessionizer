package service

type Service interface {
	CreateProject(cwd, name string)
	SelectAndOpenProject(name string)
	DeleteProject(name string)
	EditProject(name string)
}

type ServiceImpl struct{}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) CreateProject(cwd, name string) {
	panic("unimplemented")
}

func (s *ServiceImpl) SelectAndOpenProject(name string) {
	panic("unimplemented")
}

func (s *ServiceImpl) DeleteProject(name string) {
	panic("unimplemented")
}

func (s *ServiceImpl) EditProject(name string) {
	panic("unimplemented")
}
