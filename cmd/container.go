package main

import (
  "github.com/kleijnweb/git-mirror-manager/internal/git"
  "github.com/kleijnweb/git-mirror-manager/internal/http"
  "github.com/kleijnweb/git-mirror-manager/internal/manager"
  "github.com/kleijnweb/git-mirror-manager/internal/util"
)

// Container is a dead-simple DI container
type Container struct {
  config  *manager.Config
  git     git.CommandRunner
  fs      util.FileSystemUtil
  server  *http.Server
  manager *manager.Manager
}

func (c *Container) Config() *manager.Config {
  if nil == c.config {
    c.config = manager.NewConfig()
  }
  return c.config
}

func (c *Container) Git() git.CommandRunner {
  if nil == c.git {
    c.git = &git.DefaultCommandRunner{}
  }
  return c.git
}

func (c *Container) Fs() util.FileSystemUtil {
  if nil == c.fs {
    c.fs = &util.OsFileSystemUtil{}
  }
  return c.fs
}

func (c *Container) Server() *http.Server {
  if nil == c.server {
    c.server = http.NewServer(c.Manager())
  }
  return c.server
}

func (c *Container) Manager() *manager.Manager {
  if nil == c.manager {
    c.manager = manager.NewManager(
      c.Config(),
      func(uri string) (*git.Mirror, *util.ApplicationError) {
        return git.NewMirror(
          uri,
          c.Config().MirrorBaseDir,
          c.Config().MirrorUpdateInterval,
          c.Git(),
          c.Fs(),
        )
      },
      c.Git(),
      c.Fs(),
    )
  }
  return c.manager
}
