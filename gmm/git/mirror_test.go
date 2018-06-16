package git_test

import (
	"github.com/kleijnweb/git-mirror-manager/gmm"
	"github.com/kleijnweb/git-mirror-manager/gmm/git"
	"github.com/kleijnweb/git-mirror-manager/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

var uriToMirrorTestData = []struct {
	uri      string
	expected string
}{
	{"git@github.com/some/repo-Name.cmd", "some/repo-name"},
	{"git@github.com/no/cmd-suffix", "no/cmd-suffix"},
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
	{"some/repo-name", "git@github.com/some/repo-Name.cmd", false},
	{"moby/moby", "https://github.com/moby/moby", true},
	{"the-namespace/the-reponame", "https://github.com/the-namespace/the-reponame", false},
}

var updateInterval = "fauxValue"
var cronMock = &mocks.Cron{}
var gitCommandRunnerMock = &mocks.CommandRunner{}
var fsUtilMock = &mocks.FileSystemUtil{}
var updateCronFactoryStub = func(mirror *git.Mirror, interval string) (git.Cron, gmm.ApplicationError) {
	// Stubs
	cronMock.On("Start")
	cronMock.On("Stop")
	return cronMock, nil
}

func NewTestMirror(uri string, baseDir string) *git.Mirror {
	mirror, _ := git.NewMirror(
		uri,
		baseDir,
		updateInterval,
		func() *mocks.CommandRunner {
			// Stubs
			gitCommandRunnerMock.On("LsRemoteTags", mock.Anything).Return("", nil)
			gitCommandRunnerMock.On("CreateMirror", mock.Anything, mock.Anything).Return(nil)
			return gitCommandRunnerMock
		}(),
		func() *mocks.FileSystemUtil {
			// Stubs
			fsUtilMock.On("DirectoryExists", mock.Anything).Return(false)
			return fsUtilMock
		}(),
		updateCronFactoryStub,
	)

	return mirror
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
			mirror := NewTestMirror(tt.uri, baseDir)
			assert.New(t).Nil(mirror.AssertValidRemote(tt.uri))
		})
	}
}

func TestInitWillFailWhenUriIsEmpty(t *testing.T) {
	_, err := git.NewMirror(
		"",
		"/baseuri",
		"fauxvalue",
		gitCommandRunnerMock,
		fsUtilMock,
		updateCronFactoryStub,
	)
	if err == nil {
		t.Error("expected errors, got nil")
	}
	if err.Code() != gmm.ErrUser {
		t.Errorf("expected errors code %d, got %d", gmm.ErrUser, err.Code())
	}
}

func TestInitWillInitializeFields(t *testing.T) {

	mirrorBaseDir := "/some/path"
	url := "http://example.org/namespace/Name"
	mirror := NewTestMirror(url, mirrorBaseDir)

	if mirror.Name == "" {
		t.Error("mirror Name was not initialized")
	}
	if !strings.HasPrefix(mirror.Path(), mirrorBaseDir) {
		t.Errorf("expected path to be prefixed with %s, got %s", mirrorBaseDir, mirror.Path())
	}
}

func TestNewMirrorStartsCron(t *testing.T) {
	NewTestMirror("http://example.com/some/repo", "/path")
	cronMock.AssertCalled(t, "Start")
}

func TestRemoveMirrorStopsCron(t *testing.T) {
	mirror := NewTestMirror("http://example.com/some/repo", "/path")
	mirror.Destroy()
	cronMock.AssertCalled(t, "Stop")
}
