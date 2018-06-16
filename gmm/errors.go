package gmm

import (
	"errors"
	"fmt"
)

const (
	// ErrFilesystem error interacting with the file system
	ErrFilesystem = iota
	// ErrNet network error
	ErrNet = iota
	// ErrGitCommand error invoking a Git command
	ErrGitCommand = iota
	// ErrCron error scheduling a cron job
	ErrCron = iota
	// ErrUser user input error
	ErrUser = iota
	// ErrNotFound requested resource was not found
	ErrNotFound = iota
)

// ApplicationError some application error
type ApplicationError interface {
	error
	Code() int
}

// CategorizedError wraps a standard errors to add an interpretable error code
type CategorizedError struct {
	Err  error
	code int
}

// Error is an error interface method
func (e CategorizedError) Error() string {
	return fmt.Sprintf("%s [%d]", e.Err.Error(), e.code)
}

// Code returns the error code
func (e CategorizedError) Code() int {
	return e.code
}

// NewError creates a new CategorizedError from text
func NewError(text string, code int) ApplicationError {
	return NewErrorUsingError(errors.New(text), code)
}

// NewErrorUsingError creates a new CategorizedError from a standard library error
func NewErrorUsingError(error error, code int) ApplicationError {
	return &CategorizedError{error, code}
}
