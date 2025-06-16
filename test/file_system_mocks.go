package test

import (
	"io/fs"
	"os"

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
