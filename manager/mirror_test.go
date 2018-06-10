package manager

import (
	"strings"
	"testing"
)

var uriToMirrorTestData = []struct {
	uri      string
	expected string
}{
	{"git@github.com/some/repo-name.git", "some/repo-name"},
	{"git@github.com/no/git-suffix", "no/git-suffix"},
	{"https://github.com/the-namespace/the-reponame", "the-namespace/the-reponame"},
	{"https://github.com/upperCase/WillBeMadeLowerCase", "uppercase/willbemadelowercase"},
}

// TODO: remove external dependency
var validRemoteTestData = []struct {
	name  string
	uri   string
	valid bool
}{
	{"docker-library/golang", "https://github.com/docker-library/golang.git", true},
	{"some/repo-name", "git@github.com/some/repo-name.git", false},
	{"moby/moby", "https://github.com/moby/moby", true},
	{"the-namespace/the-reponame", "https://github.com/the-namespace/the-reponame", false},
}

func TestMirrorNameFromUri(t *testing.T) {
	for _, tt := range uriToMirrorTestData {
		t.Run(tt.uri, func(t *testing.T) {
			actual := mirrorNameFromURI(tt.uri)
			if actual != tt.expected {
				t.Errorf("got %q, want %q", actual, tt.expected)
			}
		})
	}
}

func TestAssertValidRemote(t *testing.T) {
	mirror := &mirror{}
	for _, tt := range validRemoteTestData {
		t.Run(tt.name, func(t *testing.T) {
			err := mirror.assertValidRemote(tt.uri)
			if err != nil && tt.valid || err == nil && !tt.valid {
				t.Errorf("expected %v result, got %v", tt.valid, !tt.valid)
			}
		})
	}
}

func TestInitWillFailWhenUriIsEmpty(t *testing.T) {
	mirror := &mirror{}
	err := mirror.init(&config{})
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.code != errUser {
		t.Errorf("expected error code %d, got %d", errUser, err.code)
	}
}

func TestInitWillInitializeFields(t *testing.T) {
	mirror := &mirror{URI: "http://example.org/namespace/name"}
	// Ignore errors for this test
	config := &config{mirrorUpdateInterval: "fauxValue", mirrorBaseDir: "/some/path"}
	mirror.init(config)
	if mirror.Name == "" {
		t.Error("mirror name was not initialized")
	}
	if mirror.UpdateInterval != config.mirrorUpdateInterval {
		t.Errorf("expected update interval %s, got %s", mirror.UpdateInterval, config.mirrorUpdateInterval)
	}
	if !strings.HasPrefix(mirror.Path, config.mirrorBaseDir) {
		t.Errorf("expected path to prefixed with %s, got %s", config.mirrorBaseDir, mirror.Path)
	}
}
