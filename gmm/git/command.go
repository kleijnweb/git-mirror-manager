package git

import (
	"github.com/kleijnweb/git-mirror-manager/gmm"
	"github.com/kleijnweb/git-mirror-manager/gmm/util"
	log "github.com/sirupsen/logrus"
	"path"
)

// CommandError represents an error executing a Git command
type CommandError interface {
	gmm.ApplicationError
}

// CommandRunner invokes the Git CLI
type CommandRunner interface {
	GetRemote(directory string) (string, CommandError)
	LsRemoteTags(uri string) (string, CommandError)
	FetchPrune(directory string) CommandError
	CreateMirror(uri string, dirPath string) CommandError
	CreateTagArchive(tag string, dirPath string) CommandError
	Exec(directory string, args ...string) (string, CommandError)
}

// DefaultCommandRunner is the default implementation of CommandRunner
type DefaultCommandRunner struct {
	Fs       util.FileSystemUtil
	Executor util.CommandExecutor
}

// GetRemote fetches the URI for the default remote at given path
func (m *DefaultCommandRunner) GetRemote(directory string) (string, CommandError) {
	return m.Exec(directory, "config", "--get", "remote.origin.url")
}

// LsRemoteTags lists the tags in a remote repository
func (m *DefaultCommandRunner) LsRemoteTags(uri string) (string, CommandError) {
	return m.Exec("", "ls-remote", "--tags", uri)
}

// FetchPrune updates a local repository with the default remote
func (m *DefaultCommandRunner) FetchPrune(directory string) CommandError {
	_, err := m.Exec(directory, "fetch", "--prune")
	return err
}

// CreateMirror creates a Git mirror on the filesystem
func (m *DefaultCommandRunner) CreateMirror(uri string, dirPath string) CommandError {
	if err := m.Fs.Mkdir(path.Dir(dirPath)); err != nil {
		return gmm.NewErrorUsingError(err, gmm.ErrFilesystem)
	}
	_, err := m.Exec("", "clone", "--mirror", "--bare", uri, dirPath)
	return err
}

// CreateTagArchive builds a ZIP file for a given tag
func (m *DefaultCommandRunner) CreateTagArchive(tag string, dirPath string) CommandError {
	if err := m.Fs.Mkdir(dirPath + "/dist/"); err != nil {
		return gmm.NewErrorUsingError(err, gmm.ErrFilesystem)
	}
	_, err := m.Exec(dirPath, "archive", tag, "-o", dirPath+"/dist/"+tag+".zip")
	return err
}

// Exec executes "git" binary commands
func (m *DefaultCommandRunner) Exec(directory string, args ...string) (string, CommandError) {
	stringOutput, err := m.Executor.Exec("git", directory, args...)

	if err != nil {
		log.Warn("Git said: " + stringOutput)
		return "", gmm.NewErrorUsingError(err, gmm.ErrFilesystem)
	}

	return stringOutput, nil
}
