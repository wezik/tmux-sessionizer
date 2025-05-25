package yaml

import (
	. "phopper/src/domain/model"
)

type YamlStorage struct{}

func NewYamlStorage() *YamlStorage {
	return &YamlStorage{}
}

func (s *YamlStorage) List() ([]*Project, error) {
	panic("unimplemented")
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
