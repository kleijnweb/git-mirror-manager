package util

import (
  "os"
  "io/ioutil"
)

type FileSystemUtil interface {
  DirectoryExists(path string) bool
  Mkdir(path string) error
  ReadDir(path string) ([]os.FileInfo, error)
}

type OsFileSystemUtil struct{}

func (u OsFileSystemUtil) DirectoryExists(path string) bool {
  if _, err := os.Stat(path); err != nil {
    return ! os.IsNotExist(err)
  }

  return true
}

func (u OsFileSystemUtil) Mkdir(path string) error {
  return os.MkdirAll(path, 700)
}

func (u OsFileSystemUtil) ReadDir(path string) ([]os.FileInfo, error) {
  return ioutil.ReadDir(path)
}
