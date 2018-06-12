package git

import (
  "github.com/kleijnweb/git-mirror-manager/internal"
  "github.com/kleijnweb/git-mirror-manager/internal/util"
  log "github.com/sirupsen/logrus"
  "path"
)

type CommandError interface {
  internal.ApplicationError
}

type CommandRunner interface {
  GetRemote(directory string) (string, CommandError)
  LsRemoteTags(uri string) (string, CommandError)
  FetchPrune(directory string) (CommandError)
  CreateMirror(uri string, dirPath string) (CommandError)
  CreateTagArchive(tag string, dirPath string) (CommandError)
  Exec(directory string, args ...string) (string, CommandError)
}

type DefaultCommandRunner struct {
  Fs       util.FileSystemUtil
  Executor util.CommandExecutor
}

func (m *DefaultCommandRunner) GetRemote(directory string) (string, CommandError) {
  return m.Exec(directory, "config", "--get", "remote.origin.url")
}

func (m *DefaultCommandRunner) LsRemoteTags(uri string) (string, CommandError) {
  return m.Exec("", "ls-remote", "--tags", uri)
}

func (m *DefaultCommandRunner) FetchPrune(directory string) (CommandError) {
  _, err := m.Exec(directory, "fetch", "--prune")
  return err
}

func (m *DefaultCommandRunner) CreateMirror(uri string, dirPath string) (CommandError) {
  if err := m.Fs.Mkdir(path.Dir(dirPath)); err != nil {
    return &internal.CategorizedError{err, internal.ErrFilesystem}
  }
  _, err := m.Exec("", "clone", "--mirror", "--bare", uri, dirPath)
  return err
}

func (m *DefaultCommandRunner) CreateTagArchive(tag string, dirPath string) (CommandError) {
  if err := m.Fs.Mkdir(dirPath + "/dist/"); err != nil {
    return &internal.CategorizedError{err, internal.ErrFilesystem}
  }
  _, err := m.Exec(dirPath, "archive", tag, "-o", dirPath+"/dist/"+tag+".zip")
  return err
}

func (m *DefaultCommandRunner) Exec(directory string, args ...string) (string, CommandError) {
  stringOutput, err := m.Executor.Exec("cmd", directory, args...)

  if err != nil {
    log.Warn("Git said: " + stringOutput)
    return "", &internal.CategorizedError{err, internal.ErrGitCommand}
  }

  return stringOutput, nil
}
