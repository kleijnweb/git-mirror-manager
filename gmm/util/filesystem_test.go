package util

import (
  "github.com/stretchr/testify/assert"
  "io/ioutil"
  "log"
  "os"
  "testing"
)

func TestOsFileSystemUtil_DirectoryExists(t *testing.T) {
  assertions := assert.New(t)
  u := OsFileSystemUtil{}
  dir, _ := ioutil.TempDir(os.TempDir(), "prefix")
  assertions.True(u.DirectoryExists(dir))
  os.Remove(dir)
  assertions.False(u.DirectoryExists(dir))
}

func TestOsFileSystemUtil_Mkdir(t *testing.T) {
  dir, _ := ioutil.TempDir(os.TempDir(), "prefix")
  defer os.Remove(dir)
  assertions := assert.New(t)
  u := OsFileSystemUtil{}
  assertions.Nil(u.Mkdir(dir+"/subdir"))
}

func TestOsFileSystemUtil_ReadDir(t *testing.T) {
  dir, _ := ioutil.TempDir(os.TempDir(), "prefix")
  defer os.Remove(dir)
  assertions := assert.New(t)
  u := OsFileSystemUtil{}

  file, err := os.OpenFile(dir +"/test.txt", os.O_WRONLY|os.O_CREATE, 0644)
  if err != nil {
    log.Fatalf("failed opening file: %s", err)
  }
  defer file.Close()

  slice, err := u.ReadDir(dir)
  assertions.Nil(err)

  assertions.Equal(slice[0].Name(), "test.txt")
}
