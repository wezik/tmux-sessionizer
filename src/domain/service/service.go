package service

import (
	"fmt"
	"phopper/src/domain/model"
)

type Service interface {
	CreateProject(cwd string, name string) (*model.Project, error)

}

type ServiceImpl struct{}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) CreateProject(cwd string, name string) (*model.Project, error) {
	fmt.Println("creating project")
	return nil, nil
}
