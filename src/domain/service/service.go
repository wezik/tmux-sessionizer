package service

import (
	"phopper/src/domain/model"
)

type Service interface {
	CreateProject(cwd, name string) (*model.Project)
	SelectAndOpenProject(name string)
}

type ServiceImpl struct{}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) CreateProject(cwd, name string) (*model.Project) {
	panic("unimplemented")
}

func (s *ServiceImpl) SelectAndOpenProject(name string) {
	panic("unimplemented")
}
