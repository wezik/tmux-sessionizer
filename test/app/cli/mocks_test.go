package cli_test

type MockService struct {
	SelectAndOpenProjectParam1 string
	SelectAndOpenProjectCalls  int

	CreateProjectParam1 string
	CreateProjectParam2 string
	CreateProjectCalls  int

	DeleteProjectParam1 string
	DeleteProjectCalls  int

	EditProjectParam1 string
	EditProjectParam2 string
	EditProjectCalls  int
}

func (s *MockService) SelectAndOpenProject(name string) {
	s.SelectAndOpenProjectParam1 = name
	s.SelectAndOpenProjectCalls++
}

func (s *MockService) CreateProject(cwd, name string) {
	s.CreateProjectParam1 = cwd
	s.CreateProjectParam2 = name
	s.CreateProjectCalls++
}

func (s *MockService) DeleteProject(name string) {
	s.DeleteProjectParam1 = name
	s.DeleteProjectCalls++
}

func (s *MockService) EditProject(name string) {
	s.EditProjectParam1 = name
	s.EditProjectCalls++
}
