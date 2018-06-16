package main

import (
	"github.com/kleijnweb/git-mirror-manager/gmm"
	"github.com/kleijnweb/git-mirror-manager/gmm/git"
	"github.com/kleijnweb/git-mirror-manager/gmm/http"
	"github.com/kleijnweb/git-mirror-manager/gmm/manager"
	"github.com/kleijnweb/git-mirror-manager/gmm/util"
)

// Container is a dead-simple DI container
type Container struct {
	config  *gmm.Config
	git     git.CommandRunner
	fs      util.FileSystemUtil
	server  *http.Server
	manager *manager.Manager
}

// Config creates and/or returns a new Config object
func (c *Container) Config() *gmm.Config {
	if nil == c.config {
		c.config = gmm.NewConfig()
	}
	return c.config
}

// Git creates and/or returns a new Git object
func (c *Container) Git() git.CommandRunner {
	if nil == c.git {
		c.git = &git.DefaultCommandRunner{}
	}
	return c.git
}

// Fs creates and/or returns a new Fs object
func (c *Container) Fs() util.FileSystemUtil {
	if nil == c.fs {
		c.fs = &util.OsFileSystemUtil{}
	}
	return c.fs
}

// Server creates and/or returns a new Server object
func (c *Container) Server() *http.Server {
	if nil == c.server {
		c.server = http.NewServer(c.Manager())
	}
	return c.server
}

// Manager creates and/or returns a new Manager object
func (c *Container) Manager() *manager.Manager {
	if nil == c.manager {
		c.manager = manager.NewManager(
			func(uri string) (*git.Mirror, gmm.ApplicationError) {
				return git.NewMirror(
					uri,
					c.Config().MirrorBaseDir,
					c.Config().MirrorUpdateInterval,
					c.Git(),
					c.Fs(),
					git.CreateUpdateCron,
				)
			},
			c.Git(),
			c.Fs(),
		)
	}
	return c.manager
}

func main() {
	container := &Container{}
	server := container.Server()
	server.Start(container.Config())
}
