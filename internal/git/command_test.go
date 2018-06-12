package git_test

import (
  "errors"
  "github.com/kleijnweb/git-mirror-manager/internal/git"
  "github.com/kleijnweb/git-mirror-manager/mocks"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "testing"
)

func factory() (*git.DefaultCommandRunner, *mocks.FileSystemUtil, *mocks.CommandExecutor) {
  mockFs := &mocks.FileSystemUtil{}
  mockExec := &mocks.CommandExecutor{}
  mockFs.On("mkdir", mock.AnythingOfType("string")).Return(nil)
  return &git.DefaultCommandRunner{Fs: mockFs, Executor: mockExec}, mockFs, mockExec
}

func TestGitExec(t *testing.T) {
  cmd, _, mockExec := factory()

  path := "/some/fauxpath"

  mockExec.On("exec", "cmd", path, "something", "--param1", "--param2").Return("", nil)

  if _, err := cmd.Exec(path, "something", "--param1", "--param2"); err != nil {
    t.Errorf("unexpected errors: %s", err)
  }

  path = "/some/other/fauxpath"

  mockExec.On("exec", "cmd", path, "something", "--param1", "--param2").Return("stderr output", errors.New("errors message"))

  if _, err := cmd.Exec(path, "something", "--param1", "--param2"); err == nil {
    t.Errorf("expected errors")
  }
}

func TestGitCreateMirror(t *testing.T) {
  cmd, _, mockExec := factory()

  uri := "https://github.com/sirupsen/logrus"
  path := "/some/fauxpath"

  mockExec.On("exec", "cmd", "", "clone", "--mirror", "--bare", uri, path).Return("", nil)

  if err := cmd.CreateMirror(uri, path); err != nil {
    t.Errorf("unexpected errors: %s", err)
  }

  path = "/some/other/fauxpath"

  mockExec.On("exec", "cmd", "", "clone", "--mirror", "--bare", uri, path).Return("stderr output", errors.New("errors message"))

  if err := git.CreateMirror(uri, path); err == nil {
    t.Errorf("expected errors")
  }
}

func TestGitCreateTagArchive(t *testing.T) {
  git, _, mockExec := factory()
  path := "/some/fauxpath"
  tag := "v1.0.0"
  mockExec.On("exec", "cmd", path, "archive", tag, "-o", path+"/dist/"+tag+".zip").Return("", nil)
  if err := git.CreateTagArchive(tag, path); err != nil {
    t.Errorf("unexpected errors: %s", err)
  }
}

func TestGitLsRemoteTags(t *testing.T) {
  git, _, mockExec := factory()
  expected := "lklk"

  uri := "https://github.com/sirupsen/logrus"
  mockExec.On("exec", "cmd", "", "ls-remote", "--tags", uri).Return(expected, nil)
  output, err := git.LsRemoteTags(uri)
  if err != nil {
    t.Errorf("unexpected errors: %s", err)
  }

  assert.New(t).Equal(expected, output, "they should be equal")
}

func TestGitFetchPrune(t *testing.T) {
  git, _, mockExec := factory()
  path := "/some/fauxpath"
  mockExec.On("exec", "cmd", path, "fetch", "--prune").Return("", nil)
  if err := git.FetchPrune(path); err != nil {
    t.Errorf("unexpected errors: %s", err)
  }
}
