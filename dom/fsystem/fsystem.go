package fsystem

import "os"

type FileSystem interface {
	MkdirAll(path string) error
	ReadDir(path string) ([]os.DirEntry, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	RemoveAll(path string) error
}
