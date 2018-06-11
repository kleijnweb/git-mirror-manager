package git

import (
  "testing"
  "github.com/stretchr/testify/mock"
  "errors"
  "github.com/stretchr/testify/assert"
)

func factory() (*DefaultCommandRunner, *mockFileSystemUtil, *mockCommandExecutor) {
  mockFs := &mockFileSystemUtil{}
  mockExec := &mockCommandExecutor{}
  mockFs.On("mkdir", mock.AnythingOfType("string")).Return(nil)
  return &DefaultCommandRunner{fs: mockFs, executor: mockExec}, mockFs, mockExec
}

func TestGitExec(t *testing.T) {
  git, _, mockExec := factory()

  path := "/some/fauxpath"

  mockExec.On("exec", "cmd", path, "something", "--param1", "--param2").Return("", nil)

  if _, err := git.exec(path, "something", "--param1", "--param2"); err != nil {
    t.Errorf("unexpected error: %s", err)
  }

  path = "/some/other/fauxpath"

  mockExec.On("exec", "cmd", path, "something", "--param1", "--param2").Return("stderr output", errors.New("error message"))

  if _, err := git.exec(path, "something", "--param1", "--param2"); err == nil {
    t.Errorf("expected error")
  }
}


func TestGitCreateMirror(t *testing.T) {
  git, _, mockExec := factory()

  uri := "https://github.com/sirupsen/logrus"
  path := "/some/fauxpath"

  mockExec.On("exec", "cmd", "", "clone", "--mirror", "--bare", uri, path).Return("", nil)

  if err := git.createMirror(uri, path); err != nil {
    t.Errorf("unexpected error: %s", err)
  }

  path = "/some/other/fauxpath"

  mockExec.On("exec", "cmd", "", "clone", "--mirror", "--bare", uri, path).Return("stderr output", errors.New("error message"))

  if err := git.createMirror(uri, path); err == nil {
    t.Errorf("expected error")
  }
}

func TestGitCreateTagArchive(t *testing.T) {
  git, _, mockExec := factory()
  path := "/some/fauxpath"
  tag := "v1.0.0"
  mockExec.On("exec", "cmd", path, "archive", tag, "-o", path+"/dist/"+tag+".zip").Return("", nil)
  if err := git.createTagArchive(tag, path); err != nil {
    t.Errorf("unexpected error: %s", err)
  }
}

func TestGitLsRemoteTags(t *testing.T) {
  git, _, mockExec := factory()
  expected := "lklk"

  uri := "https://github.com/sirupsen/logrus"
  mockExec.On("exec", "cmd", "", "ls-remote", "--tags", uri).Return(expected, nil)
  output, err := git.lsRemoteTags(uri)
  if err != nil {
    t.Errorf("unexpected error: %s", err)
  }

  assert.New(t).Equal(expected, output, "they should be equal")
}

func TestGitFetchPrune(t *testing.T) {
  git, _, mockExec := factory()
  path := "/some/fauxpath"
  mockExec.On("exec", "cmd", path, "fetch", "--prune").Return("", nil)
  if err := git.fetchPrune(path); err != nil {
    t.Errorf("unexpected error: %s", err)
  }
}

