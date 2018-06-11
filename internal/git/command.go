package git

import (
  "github.com/kleijnweb/git-mirror-manager/internal/util"
  log "github.com/sirupsen/logrus"
  "path"
)

type CommandRunner interface {
  GetRemote(directory string) (string, *util.ApplicationError)
  lsRemoteTags(uri string) (string, *util.ApplicationError)
  fetchPrune(directory string) (*util.ApplicationError)
  createMirror(uri string, dirPath string) (*util.ApplicationError)
  createTagArchive(tag string, dirPath string) (*util.ApplicationError)
  exec(directory string, args ...string) (string, *util.ApplicationError)
}

type DefaultCommandRunner struct {
  fs       util.FileSystemUtil
  executor util.CommandExecutor
}

func (m *DefaultCommandRunner) GetRemote(directory string) (string, *util.ApplicationError) {
  return m.exec(directory, "config", "--get", "remote.origin.url")
}

func (m *DefaultCommandRunner) lsRemoteTags(uri string) (string, *util.ApplicationError) {
  return m.exec("", "ls-remote", "--tags", uri)
}

func (m *DefaultCommandRunner) fetchPrune(directory string) (*util.ApplicationError) {
  _, err := m.exec(directory, "fetch", "--prune")
  return err
}

func (m *DefaultCommandRunner) createMirror(uri string, dirPath string) (*util.ApplicationError) {
  if err := m.fs.Mkdir(path.Dir(dirPath)); err != nil {
    return &util.ApplicationError{err, util.ErrFilesystem}
  }
  _, err := m.exec("", "clone", "--mirror", "--bare", uri, dirPath)
  return err
}

func (m *DefaultCommandRunner) createTagArchive(tag string, dirPath string) (*util.ApplicationError) {
  if err := m.fs.Mkdir(dirPath + "/dist/"); err != nil {
    return &util.ApplicationError{err, util.ErrFilesystem}
  }
  _, err := m.exec(dirPath, "archive", tag, "-o", dirPath+"/dist/"+tag+".zip")
  return err
}

func (m *DefaultCommandRunner) exec(directory string, args ...string) (string, *util.ApplicationError) {
  stringOutput, err := m.executor.Exec("cmd", directory, args...)

  if err != nil {
    log.Warn("Git said: " + stringOutput)
    return "", &util.ApplicationError{err, util.ErrGitCommand}
  }

  return stringOutput, nil
}
