package main

import (
  "github.com/kleijnweb/git-mirror-manager/gmm"
  "github.com/kleijnweb/git-mirror-manager/gmm/git"
  "github.com/kleijnweb/git-mirror-manager/gmm/http"
  "github.com/kleijnweb/git-mirror-manager/gmm/manager"
  "github.com/kleijnweb/git-mirror-manager/gmm/util"
  "github.com/stretchr/testify/assert"
  "testing"
)

var container = &Container{}

func TestContainer_Config(t *testing.T) {
  assert.New(t).IsType(&gmm.Config{}, container.Config())
}

func TestContainer_Fs(t *testing.T) {
  assert.New(t).IsType(&util.OsFileSystemUtil{}, container.Fs())
}

func TestContainer_Git(t *testing.T) {
  assert.New(t).IsType(&git.DefaultCommandRunner{}, container.Git())
}

func TestContainer_Manager(t *testing.T) {
  assert.New(t).IsType(&manager.Manager{}, container.Manager())
}

func TestContainer_Server(t *testing.T) {
  assert.New(t).IsType(&http.Server{}, container.Server())
}
