package yaml_test

import (
	"io/fs"
	"os"
	. "phopper/dom/model"
	. "phopper/dom/utils"
	. "phopper/infra/yaml"
	"slices"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

type MockFileSystem struct {
	MkdirAllCalls int
	MkdirAllPath  string
	MkdirAllErr   error

	ReadDirCalls  int
	ReadDirReturn []MockDirEntry
	ReadDirErr    error

	ReadFileCalls  int
	ReadFileReturn []byte
	ReadFileErr    error

	WriteFileCalls int
	WriteFilePath  string
	WriteFileBytes []byte
	WriteFileErr   error

	RemoveAllCalls int
	RemoveAllPath  string
	RemoveAllErr   error
}

func (s *MockFileSystem) MkdirAll(path string) error {
	s.MkdirAllCalls++
	s.MkdirAllPath = path
	return s.MkdirAllErr
}

func (s *MockFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	s.ReadDirCalls++
	dirEntries := make([]os.DirEntry, len(s.ReadDirReturn))
	for i, dirEntry := range s.ReadDirReturn {
		dirEntries[i] = dirEntry
	}
	return dirEntries, s.ReadDirErr
}

func (s *MockFileSystem) ReadFile(path string) ([]byte, error) {
	s.ReadFileCalls++
	return s.ReadFileReturn, s.ReadFileErr
}

func (s *MockFileSystem) WriteFile(path string, data []byte) error {
	s.WriteFileCalls++
	s.WriteFilePath = path
	s.WriteFileBytes = data
	return s.WriteFileErr
}

func (s *MockFileSystem) RemoveAll(path string) error {
	s.RemoveAllCalls++
	s.RemoveAllPath = path
	return s.RemoveAllErr
}

type MockDirEntry struct {
	isDir   bool
	name    string
	type_   os.FileMode
	info    fs.FileInfo
	infoErr error
}

func (s MockDirEntry) IsDir() bool {
	return s.isDir
}

func (s MockDirEntry) Name() string {
	return s.name
}

func (s MockDirEntry) Type() os.FileMode {
	return s.type_
}

func (s MockDirEntry) Info() (fs.FileInfo, error) {
	return s.info, s.infoErr
}

type MockConfig struct {
	GetConfigDirCalls  int
	GetConfigDirReturn string
}

func (c *MockConfig) GetConfigDir() string {
	c.GetConfigDirCalls++
	return c.GetConfigDirReturn
}

func Test_YamlStorage(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		t.Run("creates template directory", func(t *testing.T) {
			// given
			cfg := &MockConfig{}
			cfg.GetConfigDirReturn = "/foo/bar"

			fs := &MockFileSystem{}

			st := NewYamlStorage(cfg, fs)

			// when
			_, err := st.List()

			// then
			Assert(t, err == nil, "Error should be nil")
			Assert(t, fs.MkdirAllCalls == 1, "MkdirAll should be called once")
			path := fs.MkdirAllPath
			Assert(t, path == "/foo/bar/templates", "MkdirAll param should be %s is %s", "/foo/bar/templates", path)
		})

		t.Run("returns empty list when no templates exist", func(t *testing.T) {
			// given
			cfg := &MockConfig{}

			fs := &MockFileSystem{}
			fs.ReadDirReturn = []MockDirEntry{}

			st := NewYamlStorage(cfg, fs)

			// when
			projects, err := st.List()

			// then
			Assert(t, fs.ReadDirCalls == 1, "ReadDir should be called once")
			Assert(t, err == nil, "Error should be nil")
			Assert(t, len(projects) == 0, "Projects should be empty")
		})

		t.Run("returns list of projects", func(t *testing.T) {
			// given
			cfg := &MockConfig{}

			fs := &MockFileSystem{}
			fs.ReadDirReturn = []MockDirEntry{
				{name: "foo", isDir: true},
				{name: "bar", isDir: true},
				{name: "baz", isDir: true},
			}
			fs.ReadFileReturn = []byte("name: foobar\ntemplate:\n  root: /home/test\n")

			expectedProjects := []*Project{
				{ID: "foo", Name: "foobar", Template: Template{Root: "/home/test"}},
				{ID: "bar", Name: "foobar", Template: Template{Root: "/home/test"}},
				{ID: "baz", Name: "foobar", Template: Template{Root: "/home/test"}},
			}

			st := NewYamlStorage(cfg, fs)

			// when
			projects, err := st.List()

			// then
			Assert(t, err == nil, "Error should be nil")
			Assert(t, fs.ReadDirCalls == 1, "ReadDir should be called once")
			Assert(t, fs.ReadFileCalls == 3, "ReadFile should be called 3 times")

			equal := func() bool {
				if len(projects) != len(expectedProjects) {
					return false
				}

				for i, project := range expectedProjects {
					if projects[i].ID != project.ID {
						return false
					}

					if projects[i].Name != project.Name {
						return false
					}

					if projects[i].Template.Root != project.Template.Root {
						return false
					}
				}
				return true
			}()
			Assert(t, equal, "Projects should be %s is %s", expectedProjects, projects)
		})
	})

	t.Run("find", func(t *testing.T) {
		t.Run("creates template directory", func(t *testing.T) {
			// given
			cfg := &MockConfig{}
			cfg.GetConfigDirReturn = "/foo/bar"

			fs := &MockFileSystem{}

			st := NewYamlStorage(cfg, fs)

			// when
			st.Find("foobar")

			// then
			Assert(t, fs.MkdirAllCalls == 1, "MkdirAll should be called once")
			path := fs.MkdirAllPath
			Assert(t, path == "/foo/bar/templates", "MkdirAll param should be %s is %s", "/foo/bar/templates", path)
		})

		t.Run("returns ErrNotFound when no such project exists", func(t *testing.T) {
			// given
			cfg := &MockConfig{}

			fs := &MockFileSystem{}
			fs.ReadDirReturn = []MockDirEntry{}

			st := NewYamlStorage(cfg, fs)

			// when
			_, err := st.Find("foobar")

			// then
			Assert(t, err == ErrNotFound, "Error should be %s is %s", ErrNotFound, err)
		})

		t.Run("returns found project", func(t *testing.T) {
			// given
			cfg := &MockConfig{}

			fs := &MockFileSystem{}
			fs.ReadDirReturn = []MockDirEntry{
				{name: "foo", isDir: true},
			}
			fs.ReadFileReturn = []byte("name: foobar\ntemplate:\n  root: /home/test\n")
			expectedProject := &Project{ID: "foo", Name: "foobar", Template: Template{Root: "/home/test"}}

			st := NewYamlStorage(cfg, fs)

			// when
			project, err := st.Find(expectedProject.Name)

			// then
			Assert(t, err == nil, "Error should be nil")
			Assert(t, project.ID == expectedProject.ID, "Project ID should be %s is %s", expectedProject, project)
		})
	})

	t.Run("save", func(t *testing.T) {
		t.Run("creates template directory", func(t *testing.T) {
			// given
			cfg := &MockConfig{}
			cfg.GetConfigDirReturn = "/foo/bar"

			fs := &MockFileSystem{}

			st := NewYamlStorage(cfg, fs)

			// when
			err := st.Save(&Project{ID: "foo"})

			// then
			Assert(t, err == nil, "Error should be nil is %s", err)
			Assert(t, fs.MkdirAllCalls == 1, "MkdirAll should be called once")
			path := fs.MkdirAllPath
			Assert(t, path == "/foo/bar/templates/foo", "MkdirAll param should be %s is %s", "/foo/bar/templates/foo", path)
		})

		t.Run("generates UUID for new project", func(t *testing.T) {
			// given
			cfg := &MockConfig{}

			fs := &MockFileSystem{}

			st := NewYamlStorage(cfg, fs)

			// when
			err := st.Save(&Project{})

			// then
			Assert(t, err == nil, "Error should be nil is %s", err)

			// TODO: think about it, probably save should just return the saved project
			param := fs.MkdirAllPath
			strings := strings.Split(param, "/")
			Assert(t, len(strings) == 2, "Param should be %s is %s", "templates/<uuid>", param)
			Assert(t, strings[1] != "", "UUID should not be empty")
		})

		t.Run("writes template to file", func(t *testing.T) {
			// given
			cfg := &MockConfig{}

			fs := &MockFileSystem{}

			project := &Project{ID: "foo", Name: "foobar", Template: Template{Root: "/home/test"}}

			expectedBytes, err := yaml.Marshal(project)
			Assert(t, err == nil, "Error should be nil")

			expectedPath := "templates/foo/template.yaml"

			st := NewYamlStorage(cfg, fs)

			// when
			err = st.Save(project)

			// then
			Assert(t, err == nil, "Error should be nil")
			Assert(t, fs.WriteFileCalls == 1, "WriteFile should be called once")
			Assert(t, fs.WriteFilePath == expectedPath, "WriteFile path param should be %s is %s", expectedPath, fs.WriteFilePath)
			Assert(t, slices.Equal(expectedBytes, fs.WriteFileBytes), "WriteFile data param should be %s is %s", expectedBytes, fs.WriteFileBytes)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("deletes template directory", func(t *testing.T) {
			// given
			cfg := &MockConfig{}
			cfg.GetConfigDirReturn = "/foo/bar"

			fs := &MockFileSystem{}

			st := NewYamlStorage(cfg, fs)

			// when
			err := st.Delete("foo")

			// then
			Assert(t, err == nil, "Error should be nil is %s", err)
			Assert(t, fs.RemoveAllCalls == 1, "RemoveAll should be called once")
			path := fs.RemoveAllPath
			Assert(t, path == "/foo/bar/templates/foo", "RemoveAll param should be %s is %s", "/foo/bar/templates/foo", path)
		})
	})

}
