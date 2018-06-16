package util_test

import (
	"github.com/kleijnweb/git-mirror-manager/gmm/util"
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestExec(t *testing.T) {
  command := &util.OsCommandExecutor{}
  output, err := command.Exec("who", "/", "-u")
  if err != nil {
    t.Errorf("exec errored: %s", output)
  }
}

func TestExecFailure(t *testing.T) {
  command := &util.OsCommandExecutor{}
  output, err := command.Exec("this-command-does-not-exist", "/", "-u")

  assertions := assert.New(t)
  assertions.Error(err)
  assertions.Equal("", output)
}
