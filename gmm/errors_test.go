package gmm

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError("TestNewError", 12345)
	assertions := assert.New(t)
	assertions.Equal(err.Code(), 12345)
	assertions.Equal(err.Error(), "TestNewError [12345]")
}

func TestNewErrorUsingError(t *testing.T) {
	err := NewErrorUsingError(errors.New("TestNewErrorUsingError"), 6789)
	assertions := assert.New(t)
	assertions.Equal(err.Code(), 6789)
	assertions.Equal(err.Error(), "TestNewErrorUsingError [6789]")
}
