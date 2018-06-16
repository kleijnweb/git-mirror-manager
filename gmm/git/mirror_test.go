package git_test

import (
  "github.com/kleijnweb/git-mirror-manager/mocks"
  "strings"
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/kleijnweb/git-mirror-manager/gmm"
  "github.com/kleijnweb/git-mirror-manager/gmm/git"
)

var uriToMirrorTestData = []struct {
  uri      string
  expected string
}{
  {"cmd@github.com/some/repo-Name.cmd", "some/repo-Name"},
  {"cmd@github.com/no/cmd-suffix", "no/cmd-suffix"},
  {"https://github.com/the-namespace/the-reponame", "the-namespace/the-reponame"},
  {"https://github.com/upperCase/WillBeMadeLowerCase", "uppercase/willbemadelowercase"},
}

// TODO: remove external dependency
var validRemoteTestData = []struct {
  name  string
  uri   string
  valid bool
}{
  {"docker-library/golang", "https://github.com/docker-library/golang.cmd", true},
  {"some/repo-Name", "cmd@github.com/some/repo-Name.cmd", false},
  {"moby/moby", "https://github.com/moby/moby", true},
  {"the-namespace/the-reponame", "https://github.com/the-namespace/the-reponame", false},
}

func TestMirrorNameFromUri(t *testing.T) {
  for _, tt := range uriToMirrorTestData {
    t.Run(tt.uri, func(t *testing.T) {
      actual := git.MirrorNameFromURI(tt.uri)
      assert.New(t).Equal(tt.expected, actual)
    })
  }
}

func TestAssertValidRemote(t *testing.T) {
  for _, tt := range validRemoteTestData {
    t.Run(tt.name, func(t *testing.T) {
      baseDir := "/some/path"
      path := baseDir + "/" + tt.name
      mirror, _ := git.NewMirror(
        tt.uri,
        baseDir,
        "fauxValue",
        func() *mocks.CommandRunner {
          mock := &mocks.CommandRunner{}
          mock.On("lsRemoteTags", tt.uri).Return("", nil)
          mock.On("createMirror", tt.uri, path).Return(nil)
          return mock
        }(),
        func() *mocks.FileSystemUtil {
          mock := &mocks.FileSystemUtil{}
          mock.On("directoryExists", path).Return(false)
          return mock
        }(),
      )
      assert.New(t).Nil(mirror.AssertValidRemote(tt.uri))
    })
  }
}

func TestInitWillFailWhenUriIsEmpty(t *testing.T) {
  _, err := git.NewMirror(
    "",
    "/baseuri",
    "fauxvalue",
    &mocks.CommandRunner{},
    &mocks.FileSystemUtil{},
  )
  if err == nil {
    t.Error("expected errors, got nil")
  }
  if err.Code() != gmm.ErrUser {
    t.Errorf("expected errors code %d, got %d", gmm.ErrUser, err.Code())
  }
}

func TestInitWillInitializeFields(t *testing.T) {

  mirrorUpdateInterval := "fauxValue"
  mirrorBaseDir := "/some/path"

  url := "http://example.org/namespace/Name"
  path := "/some/path/namespace/Name"
  mirror, _ := git.NewMirror(
    url,
    mirrorBaseDir,
    mirrorUpdateInterval,
    func() *mocks.CommandRunner {
      mock := &mocks.CommandRunner{}
      mock.On("lsRemoteTags", url).Return("", nil)
      mock.On("createMirror", url, path).Return(nil)
      return mock
    }(),
    func() *mocks.FileSystemUtil {
      mock := &mocks.FileSystemUtil{}
      mock.On("directoryExists", path).Return(false)
      return mock
    }(),
  )

  if mirror.Name == "" {
    t.Error("mirror Name was not initialized")
  }
  if mirror.UpdateInterval() != mirrorUpdateInterval {
    t.Errorf("expected Update interval %s, got %s", mirror.UpdateInterval(), mirrorUpdateInterval)
  }
  if !strings.HasPrefix(mirror.Path(), mirrorBaseDir) {
    t.Errorf("expected path to be prefixed with %s, got %s", mirrorBaseDir, mirror.Path())
  }
}
