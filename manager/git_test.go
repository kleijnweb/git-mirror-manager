package manager

import (
	"testing"
	"io/ioutil"
	"os"
	"errors"
)

var tMkTmp = func(t *testing.T) string {
	tmpDirPathName, ioutilErr := ioutil.TempDir("", "gmm-git-test")
	if ioutilErr != nil {
		t.Errorf("Failed to create temporary test directory")
	}
	return tmpDirPathName
}

var gitExecPassingTestData = []struct {
	name string
	args []string
}{
	{"git init", []string{"init"}},
	{"git ls-remote", []string{"ls-remote", "https://github.com/sirupsen/logrus"}},
	{"git clone", []string{"clone", "--mirror", "https://github.com/sirupsen/logrus"}},
}

var gitExecFailingTestData = []struct {
	name string
	args []string
}{
	{"git finit", []string{"finit"}},
	{"githiiiiiib123.com", []string{"ls-remote", "--tags", "https://githiiiiiib123.com/sirupsen/logrus"}},
}

var withFailingGitDirMaker = func(t *testing.T, fn func(t *testing.T)([]byte, error)) {
	origDirMaker := gitDirMaker

	gitDirMaker = func(path string, perm os.FileMode) error {
		return errors.New("faux error")
	}

	output, err := fn(t)

	if err == nil {
		t.Errorf("Expected error, got: %s", string(output))
	}

	if err.Error() != "faux error [0]" {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	// Restore gitDirMaker
	gitDirMaker = origDirMaker
}

func TestGitExec(t *testing.T) {
	for _, tt := range gitExecPassingTestData {
		t.Run(tt.name, func(t *testing.T) {
			output, err := gitExec(tMkTmp(t), tt.args...)
			if err != nil {
				t.Errorf("gitExec errored: %s", string(output))
			}
		})
	}

	for _, tt := range gitExecFailingTestData {
		t.Run(tt.name, func(t *testing.T) {
			output, err := gitExec(tMkTmp(t), tt.args...)
			if err == nil {
				t.Errorf("expected gitExec error, got: %s", string(output))
			}
		})
	}
}

func TestGitCreateMirror(t *testing.T) {
	withFailingGitDirMaker(t, func(t *testing.T)([]byte, error) {
		return gitCreateMirror("https://github.com/sirupsen/logrus", tMkTmp(t))
	})

	output, err := gitCreateMirror("https://github.com/sirupsen/logrus", tMkTmp(t))
	if err != nil {
		t.Errorf("gitCreateMirror errored: %s", string(output))
	}
}

func TestGitCreateTagArchive(t *testing.T) {
	withFailingGitDirMaker(t, func(t *testing.T)([]byte, error) {
		tmpDirName := tMkTmp(t)
		gitCreateMirror("https://github.com/sirupsen/logrus", tmpDirName)
		return gitCreateTagArchive("v1.0.0", tmpDirName)
	})

	tmpDirName := tMkTmp(t)
	output, _ := gitCreateMirror("https://github.com/sirupsen/logrus", tmpDirName)
	output, err := gitCreateTagArchive("v1.0.0", tmpDirName)

	if err != nil {
		t.Errorf("gitCreateTagArchive errored: %s", string(output))
	}
}

func TestGitLsRemoteTags(t *testing.T) {
	output, err := gitLsRemoteTags("https://github.com/sirupsen/logrus")
	if err != nil {
		t.Errorf("gitLsRemoteTags errored: %s", string(output))
	}
}

func TestGitFetchPrune(t *testing.T) {
	tmpDirName := tMkTmp(t)
	output, _ := gitCreateMirror("https://github.com/sirupsen/logrus", tmpDirName)
	output, err := gitFetchPrune(tmpDirName)

	if err != nil {
		t.Errorf("gitFetchPrune errored: %s", string(output))
	}
}
