package git_test

import (
	"errors"
	"github.com/kleijnweb/git-mirror-manager/gmm/git"
	"github.com/kleijnweb/git-mirror-manager/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func factory() (*git.DefaultCommandRunner, *mocks.FileSystemUtil, *mocks.CommandExecutor) {
	mockFs := &mocks.FileSystemUtil{}
	mockExec := &mocks.CommandExecutor{}
	mockFs.On("Mkdir", mock.AnythingOfType("string")).Return(nil)
	return &git.DefaultCommandRunner{Fs: mockFs, Executor: mockExec}, mockFs, mockExec
}

func TestGitExec(t *testing.T) {
	cmd, _, mockExec := factory()

	path := "/some/fauxpath"

	mockExec.On("Exec", "git", path, "something", "--param1", "--param2").Return("", nil)

	if _, err := cmd.Exec(path, "something", "--param1", "--param2"); err != nil {
		t.Errorf("unexpected errors: %s", err)
	}

	path = "/some/other/fauxpath"

	mockExec.On("Exec", "git", path, "something", "--param1", "--param2").Return("stderr output", errors.New("errors message"))

	if _, err := cmd.Exec(path, "something", "--param1", "--param2"); err == nil {
		t.Errorf("expected errors")
	}
}

func TestGitCreateMirror(t *testing.T) {
	cmd, _, mockExec := factory()

	uri := "https://github.com/sirupsen/logrus"
	path := "/some/fauxpath"

	mockExec.On("Exec", "git", "", "clone", "--mirror", "--bare", uri, path).Return("", nil)

	if err := cmd.CreateMirror(uri, path); err != nil {
		t.Errorf("unexpected errors: %s", err)
	}

	path = "/some/other/fauxpath"

	mockExec.On("Exec", "git", "", "clone", "--mirror", "--bare", uri, path).Return("stderr output", errors.New("errors message"))

	if err := cmd.CreateMirror(uri, path); err == nil {
		t.Errorf("expected errors")
	}
}

func TestGitCreateTagArchive(t *testing.T) {
	cmd, _, mockExec := factory()
	path := "/some/fauxpath"
	tag := "v1.0.0"
	mockExec.On("Exec", "git", path, "archive", tag, "-o", path+"/dist/"+tag+".zip").Return("", nil)
	if err := cmd.CreateTagArchive(tag, path); err != nil {
		t.Errorf("unexpected errors: %s", err)
	}
}

func TestGitLsRemoteTags(t *testing.T) {
  cmd, _, mockExec := factory()
  expected := "lklk"

  uri := "https://github.com/sirupsen/logrus"
  mockExec.On("Exec", "git", "", "ls-remote", "--tags", uri).Return(expected, nil)
  output, err := cmd.LsRemoteTags(uri)
  if err != nil {
    t.Errorf("unexpected errors: %s", err)
  }

  assert.New(t).Equal(expected, output, "they should be equal")
}

func TestGitGetRemote(t *testing.T) {
  cmd, _, mockExec := factory()
  expected := "lklk"

  directory := "/some/path"
  mockExec.On("Exec", "git", directory, "config", "--get", "remote.origin.url").Return(expected, nil)
  output, err := cmd.GetRemote(directory)
  if err != nil {
    t.Errorf("unexpected errors: %s", err)
  }

  assert.New(t).Equal(expected, output, "they should be equal")
}

func TestGitFetchPrune(t *testing.T) {
	cmd, _, mockExec := factory()
	path := "/some/fauxpath"
	mockExec.On("Exec", "git", path, "fetch", "--prune").Return("", nil)
	if err := cmd.FetchPrune(path); err != nil {
		t.Errorf("unexpected errors: %s", err)
	}
}
