package manager

import (
	"strings"
	"os/exec"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

var gitDirMaker = os.MkdirAll

// gitExec invokes the git binary using CLI and returns STDOUT and STDERR as a byte slice.
func gitExec(directory string, args ...string) ([]byte, *Error) {
	cmd := exec.Command("git", args...)
	if directory != "" {
		cmd.Dir = directory
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Warn("Git said: " + strings.TrimSpace(string(output)))
		return nil, &Error{err, errGitCommand}
	}
	return output, nil
}

func gitFetchPrune(directory string) ([]byte, *Error) {
	return gitExec(directory, "fetch", "--prune")
}

func gitLsRemoteTags(uri string) ([]byte, *Error) {
	return gitExec("", "ls-remote", "--tags", uri)
}

func gitCreateMirror(uri string, dirPath string) ([]byte, *Error) {
	if err := gitDirMaker(path.Dir(dirPath), 0700); err != nil {
		return nil, &Error{err, errFilesystem}
	}
	return gitExec("", "clone", "--mirror", "--bare", uri, dirPath)
}

func gitCreateTagArchive(tag string, dirPath string) ([]byte, *Error) {
	if err := gitDirMaker(dirPath+"/dist/", 0700); err != nil {
		return nil, &Error{err, errFilesystem}
	}
	return gitExec(dirPath, "archive", tag, "-o", dirPath+"/dist/"+tag+".zip")
}
