package gmm

import (
	"errors"
	"fmt"
)

const (
	ErrFilesystem = iota // There was some errors interacting with the file system
	ErrNet        = iota // There was some network errors
	ErrGitCommand = iota // There was some errors invoking a Git command
	ErrCron       = iota // There was some errors scheduling a cron job
	ErrUser       = iota // There was some errors with user input
	ErrNotFound   = iota // The requested resource was not found
)

type ApplicationError interface {
  error
  Code() int
}

// CategorizedError wraps a standard errors to add an interpretable error code
type CategorizedError struct {
	Err  error
	code int
}

func (e CategorizedError) Error() string {
  return fmt.Sprintf("%s [%d]", e.Err.Error(), e.code)
}

func (e CategorizedError) Code() int {
  return e.code
}

func NewError(text string, code int) ApplicationError {
  return NewErrorUsingError(errors.New(text), code)
}

func NewErrorUsingError(error error, code int) ApplicationError {
  return &CategorizedError{error, code}
}
