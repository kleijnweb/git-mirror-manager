package util_test

import (
	"github.com/kleijnweb/git-mirror-manager/gmm/util"
	"testing"
)

func TestExec(t *testing.T) {
	command := &util.OsCommandExecutor{}
	output, err := command.Exec("who", "/", "-u")
	if err != nil {
		t.Errorf("exec errored: %s", output)
	}
}
