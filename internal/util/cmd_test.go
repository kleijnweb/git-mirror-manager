package util

import (
  "testing"
)

func TestExec(t *testing.T) {
  command := &OsCommandExecutor{}
  output, err := command.Exec("who", "/", "-u")
  if err != nil {
    t.Errorf("exec errored: %s", output)
  }
}
