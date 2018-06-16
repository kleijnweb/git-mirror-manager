package util

import (
	"io/ioutil"
	"os"
)

// FileSystemUtil wraps some common filesystem operations
type FileSystemUtil interface {
	DirectoryExists(path string) bool
	Mkdir(path string) error
	ReadDir(path string) ([]os.FileInfo, error)
}

// OsFileSystemUtil delegates to standard librarys functions
type OsFileSystemUtil struct{}

// DirectoryExists checks if a directory exists
func (u OsFileSystemUtil) DirectoryExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

// Mkdir creates a directory with permissions 700
func (u OsFileSystemUtil) Mkdir(path string) error {
	return os.MkdirAll(path, 700)
}

// ReadDir returns a slice of os.FileInfo describing directory contents
func (u OsFileSystemUtil) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}
