package service_test

import . "phopper/src/domain/model"

type MockSelector struct {
	SelectFromParam1 []string
	SelectFromParam2 string
	SelectFromCalls  int
	SelectFromReturn string
	SelectFromErr    error
}

func (s *MockSelector) SelectFrom(items []string, prompt string) (string, error) {
	s.SelectFromParam1 = items
	s.SelectFromParam2 = prompt
	s.SelectFromCalls++
	return s.SelectFromReturn, s.SelectFromErr
}

type MockMultiplexer struct {
	AttachProjectParam1 *Project
	AttachProjectCalls  int
}

func (s *MockMultiplexer) AttachProject(p *Project) error {
	s.AttachProjectParam1 = p
	s.AttachProjectCalls++
	return nil
}

type MockStorage struct {
	ListCalls  int
	ListReturn []*Project
	ListErr    error

	FindParam1 string
	FindCalls  int
	FindReturn *Project
	FindErr    error

	SaveParam1 *Project
	SaveCalls  int

	DeleteParam1 string
	DeleteCalls  int
}

func (s *MockStorage) List() ([]*Project, error) {
	s.ListCalls++
	return s.ListReturn, s.ListErr
}

func (s *MockStorage) Find(name string) (*Project, error) {
	s.FindParam1 = name
	s.FindCalls++
	return s.FindReturn, s.FindErr
}

func (s *MockStorage) Save(t *Project) error {
	s.SaveParam1 = t
	s.SaveCalls++
	return nil
}

func (s *MockStorage) Delete(uuid string) error {
	s.DeleteParam1 = uuid
	s.DeleteCalls++
	return nil
}
