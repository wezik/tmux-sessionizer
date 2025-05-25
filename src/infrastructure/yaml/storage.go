package yaml

import (
	"phopper/src/app/config"
	. "phopper/src/domain/model"
)

type YamlStorage struct{
	config *config.Config
}

func NewYamlStorage(config *config.Config) *YamlStorage {
	return &YamlStorage{config: config}
}

func (s *YamlStorage) List() ([]*Project, error) {
	p1 := &Project{Name: "foo"}
	p2 := &Project{Name: "bar"}
	p3 := &Project{Name: "baz"}
	return []*Project{p1, p2, p3}, nil
}

func (s *YamlStorage) Find(name string) (*Project, error) {
	panic("unimplemented")
}

func (s *YamlStorage) Save(t *Project) error {
	panic("unimplemented")
}

func (s *YamlStorage) Delete(uuid string) error {
	panic("unimplemented")
}
