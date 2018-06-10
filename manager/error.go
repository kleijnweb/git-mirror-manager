package manager

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	errFilesystem = iota // There was some error interacting with the file system
	errNet        = iota // There was some network error
	errGitCommand = iota // There was some error invoking a Git command
	errCron       = iota // There was some error scheduling a cron job
	errUser       = iota // There was some error with user input
	errNotFound   = iota // The requested resource was not found
)

type Error struct {
	err  error
	code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s [%d]", e.err.Error(), e.code)
}

func (e *Error) Log() {
	log.Error(e)
}

func newError(text string, code int) *Error {
	return &Error{errors.New(text), code}
}
