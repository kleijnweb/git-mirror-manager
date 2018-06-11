package util

import (
	"errors"
	"fmt"
)

const (
	ErrFilesystem = iota // There was some error interacting with the file system
	ErrNet        = iota // There was some network error
	ErrGitCommand = iota // There was some error invoking a Git command
	ErrCron       = iota // There was some error scheduling a cron job
	ErrUser       = iota // There was some error with user input
	ErrNotFound   = iota // The requested resource was not found
)

// ApplicationError wraps a standard error to add an interpretable error Code
type ApplicationError struct {
	Err  error
	Code int
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("%s [%d]", e.Err.Error(), e.Code)
}

func NewError(text string, code int) *ApplicationError {
	return &ApplicationError{errors.New(text), code}
}
