package util

import (
  "os/exec"
  "strings"
)

type CommandExecutor interface {
  Exec(name string, directory string, args ...string) (string, error)
}

type OsCommandExecutor struct{}

// Exec invokes the a binary using CLI and returns STDOUT and STDERR as a string.
func (m *OsCommandExecutor) Exec(name string, directory string, args ...string) (string, error) {
  cmd := exec.Command(name, args...)
  if directory != "" {
    cmd.Dir = directory
  }

  output, err := cmd.CombinedOutput()
  stringOutput := strings.TrimSpace(string(output))

  if err != nil {
    return "", err
  }

  return stringOutput, nil
}
