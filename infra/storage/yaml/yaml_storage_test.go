package yaml_test

import (
	"io/fs"
	"os"
	"testing"

	. "thop/dom/model"
	. "thop/infra/yaml"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileSystem struct {
	mock.Mock
}

func (s *MockFileSystem) MkdirAll(path string) error {
	args := s.Called(path)
	return args.Error(0)
}

func (s *MockFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	args := s.Called(path)
	return args.Get(0).([]os.DirEntry), args.Error(1)
}

func (s *MockFileSystem) ReadFile(path string) ([]byte, error) {
	args := s.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *MockFileSystem) WriteFile(path string, data []byte) error {
	args := s.Called(path, data)
	return args.Error(0)
}

func (s *MockFileSystem) RemoveAll(path string) error {
	args := s.Called(path)
	return args.Error(0)
}

type MockDirEntry struct {
	mock.Mock
}

func (s *MockDirEntry) IsDir() bool {
	args := s.Called()
	return args.Bool(0)
}

func (s *MockDirEntry) Name() string {
	args := s.Called()
	return args.String(0)
}

func (s *MockDirEntry) Type() os.FileMode {
	args := s.Called()
	return args.Get(0).(os.FileMode)
}

func (s *MockDirEntry) Info() (fs.FileInfo, error) {
	args := s.Called()
	return args.Get(0).(fs.FileInfo), args.Error(1)
}

type MockConfig struct {
	mock.Mock
}

func (c *MockConfig) GetConfigDir() string {
	args := c.Called()
	return args.String(0)
}

func (c *MockConfig) GetEditor() string {
	args := c.Called()
	return args.String(0)
}

func Test_List(t *testing.T) {
	cfg := new(MockConfig)
	cfg.On("GetConfigDir").Return("/foo/bar")

	t.Run("returns empty list when no templates exist", func(t *testing.T) {
		// given
		fs := new(MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{}, nil).Once()

		st := NewYamlStorage(cfg, fs)

		// when
		projects, err := st.List()

		// then
		assert.Nil(t, err)
		assert.Empty(t, projects)
		fs.AssertExpectations(t)
	})

	t.Run("returns list of projects", func(t *testing.T) {
		// given
		fs := new(MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()

		dir1 := new(MockDirEntry)
		dir1.On("IsDir").Return(true).Once()
		dir1.On("Name").Return("foo").Once()

		dir2 := new(MockDirEntry)
		dir2.On("IsDir").Return(true).Once()
		dir2.On("Name").Return("bar").Once()

		dir3 := new(MockDirEntry)
		dir3.On("IsDir").Return(true).Once()
		dir3.On("Name").Return("baz").Once()

		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir1, dir2, dir3}, nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/foo/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/bar/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/baz/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()

		expectedProjects := []*Project{
			{ID: "foo", Name: "foobar", Template: &Template{Root: "/home/test"}},
			{ID: "bar", Name: "foobar", Template: &Template{Root: "/home/test"}},
			{ID: "baz", Name: "foobar", Template: &Template{Root: "/home/test"}},
		}

		st := NewYamlStorage(cfg, fs)

		// when
		projects, err := st.List()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedProjects, projects)
		fs.AssertExpectations(t)
	})

	t.Run("gracefully handles non directory files", func(t *testing.T) {
		// given
		dir := new(MockDirEntry)
		dir.On("IsDir").Return(false).Once()

		fs := new(MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir}, nil).Once()

		st := NewYamlStorage(cfg, fs)

		// when
		projects, err := st.List()

		// then
		assert.Nil(t, err)
		assert.Empty(t, projects)
		fs.AssertExpectations(t)
	})

	t.Run("gracefully handles error while reading template file", func(t *testing.T) {
		// given
		dir := new(MockDirEntry)
		dir.On("IsDir").Return(true).Once()
		dir.On("Name").Return("foo").Once()

		fs := new(MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir}, nil).Once()

		bytes := []byte("  name: invalid:format: \ntemplate:\n  root: /home/test\n")
		fs.On("ReadFile", "/foo/bar/templates/foo/template.yaml").Return(bytes, nil).Once()

		st := NewYamlStorage(cfg, fs)

		// when
		projects, err := st.List()

		// then
		assert.Nil(t, err)
		assert.Empty(t, projects)
		fs.AssertExpectations(t)
	})
}

func Test_Find(t *testing.T) {
	cfg := new(MockConfig)
	cfg.On("GetConfigDir").Return("/foo/bar")

	t.Run("returns error when project does not exist", func(t *testing.T) {
		// given
		fs := new(MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{}, nil).Once()

		st := NewYamlStorage(cfg, fs)

		// when
		_, err := st.Find("foobar")

		// then
		assert.Equal(t, ErrNotFound, err)
		fs.AssertExpectations(t)
	})

	t.Run("returns found project", func(t *testing.T) {
		// given
		fs := new(MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()

		dir := new(MockDirEntry)
		dir.On("IsDir").Return(true).Once()
		dir.On("Name").Return("foo").Once()

		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir}, nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/foo/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		expectedProject := &Project{ID: "foo", Name: "foobar", Template: &Template{Root: "/home/test"}}

		st := NewYamlStorage(cfg, fs)

		// when
		project, err := st.Find(expectedProject.Name)

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedProject, project)
		fs.AssertExpectations(t)
	})
}

func Test_Save(t *testing.T) {
	t.Run("saves project with a template file", func(t *testing.T) {
		// given
		cfg := new(MockConfig)
		cfg.On("GetConfigDir").Return("/foo/bar")

		fs := new(MockFileSystem)
		var path string
		fs.On("MkdirAll", mock.Anything).Run(func(args mock.Arguments) {
			path = args.Get(0).(string)
		}).Return(nil).Once()

		var templatePath string
		fs.On("WriteFile", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			templatePath = args.Get(0).(string)
		}).Return(nil).Once()

		st := NewYamlStorage(cfg, fs)
		project := &Project{Name: "foo"}

		// when
		err := st.Save(project)

		// then
		assert.Nil(t, err)

		assert.NotEmpty(t, project.ID)
		assert.Equal(t, "/foo/bar/templates/"+project.ID, path)
		assert.Equal(t, "/foo/bar/templates/"+project.ID+"/template.yaml", templatePath)
	})

	t.Run("keeps UUID if it's already set", func(t *testing.T) {
		// given
		cfg := new(MockConfig)
		cfg.On("GetConfigDir").Return("/foo/bar")

		fs := new(MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates/foobar").Return(nil).Once()

		fs.On("WriteFile", "/foo/bar/templates/foobar/template.yaml", mock.Anything).Return(nil).Once()

		st := NewYamlStorage(cfg, fs)
		project := &Project{ID: "foobar", Name: "foo"}

		// when
		err := st.Save(project)

		// then
		assert.Nil(t, err)
		fs.AssertExpectations(t)
	})
}

func Test_Delete(t *testing.T) {
	cfg := new(MockConfig)
	cfg.On("GetConfigDir").Return("/foo/bar")

	t.Run("deletes template directory", func(t *testing.T) {
		// given
		fs := new(MockFileSystem)
		fs.On("RemoveAll", "/foo/bar/templates/foo").Return(nil).Once()

		st := NewYamlStorage(cfg, fs)

		// when
		err := st.Delete("foo")

		// then
		assert.Nil(t, err)
		fs.AssertExpectations(t)
	})
}

func Test_PrepareTemplateFile(t *testing.T) {
	cfg := new(MockConfig)
	cfg.On("GetConfigDir").Return("/foo/bar")

	t.Run("returns path to template file", func(t *testing.T) {
		// given
		fs := new(MockFileSystem)
		st := NewYamlStorage(cfg, fs)

		project := &Project{ID: "foo", Name: "foobar", Template: &Template{Root: "/home/test"}}

		// when
		path, err := st.PrepareTemplateFile(project)

		// then
		assert.Nil(t, err)
		assert.Equal(t, "/foo/bar/templates/foo/template.yaml", path)
		fs.AssertExpectations(t)
	})
}
