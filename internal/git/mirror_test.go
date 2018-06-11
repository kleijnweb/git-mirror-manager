package git

import (
  "strings"
  "testing"
  "github.com/stretchr/testify/assert"
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
      actual := mirrorNameFromURI(tt.uri)
      assert.New(t).Equal(tt.expected, actual)
    })
  }
}

func TestAssertValidRemote(t *testing.T) {
  for _, tt := range validRemoteTestData {
    t.Run(tt.name, func(t *testing.T) {
      baseDir := "/some/path"
      path := baseDir+"/"+tt.name
      mirror, _ := NewMirror(
        &Config{mirrorUpdateInterval: "fauxValue", mirrorBaseDir: baseDir},
        tt.uri,
        func() *mockGitCommandRunner {
          mock := &mockGitCommandRunner{}
          mock.On("lsRemoteTags", tt.uri).Return("", nil)
          mock.On("createMirror", tt.uri, path).Return(nil)
          return mock
        }(),
        func() *mockFileSystemUtil {
          mock := &mockFileSystemUtil{}
          mock.On("directoryExists", path).Return(false)
          return mock
        }(),
      )
      assert.New(t).Nil(mirror.assertValidRemote(tt.uri))
    })
  }
}

func TestInitWillFailWhenUriIsEmpty(t *testing.T) {
  _, err := NewMirror(
    &Config{},
    "",
    &mockGitCommandRunner{},
    &mockFileSystemUtil{},
  )
  if err == nil {
    t.Error("expected error, got nil")
  }
  if err.code != errUser {
    t.Errorf("expected error code %d, got %d", errUser, err.code)
  }
}

func TestInitWillInitializeFields(t *testing.T) {
  config := &Config{mirrorUpdateInterval: "fauxValue", mirrorBaseDir: "/some/path"}
  url := "http://example.org/namespace/Name"
  path := "/some/path/namespace/Name"
  mirror, _ := NewMirror(
    config,
    url,
    func() *mockGitCommandRunner {
      mock := &mockGitCommandRunner{}
      mock.On("lsRemoteTags", url).Return("", nil)
      mock.On("createMirror", url, path).Return(nil)
      return mock
    }(),
    func() *mockFileSystemUtil {
      mock := &mockFileSystemUtil{}
      mock.On("directoryExists", path).Return(false)
      return mock
    }(),
  )

  if mirror.Name == "" {
    t.Error("mirror Name was not initialized")
  }
  if mirror.updateInterval != config.mirrorUpdateInterval {
    t.Errorf("expected Update interval %s, got %s", mirror.updateInterval, config.mirrorUpdateInterval)
  }
  if !strings.HasPrefix(mirror.path, config.mirrorBaseDir) {
    t.Errorf("expected path to prefixed with %s, got %s", config.mirrorBaseDir, mirror.path)
  }
}
