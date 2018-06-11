package manager

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func withMirrorFactory(f func(config *Config, uri string) (*mirror, *ApplicationError), fn func(t *testing.T), t *testing.T) {
  orig := mirrorFactory
  mirrorFactory = f
  fn(t)
  mirrorFactory = orig
}

func TestNewManagerFailsIfBaseDirDoesNotExist(t *testing.T) {
  withMirrorFactory(
    func(config *Config, uri string) (*mirror, *ApplicationError) {
      return &mirror{}, nil
    },
    func(t *testing.T) {
      _, err := NewManager(
        &Config{MirrorBaseDir: "/this/directory/should/not/exist"},
        func() gitCommandRunner {
          return &mockGitCommandRunner{}
        }(),
      )

      a := assert.New(t)
      a.Error(err)
      a.Equal(errFilesystem, err.code)
    },
    t,
  )
}

func TestMapOperations(t *testing.T) {
  withMirrorFactory(
    func(config *Config, uri string) (*mirror, *ApplicationError) {
      return &mirror{}, nil
    },
    func(t *testing.T) {
      _, err := NewManager(
        &Config{MirrorBaseDir: "/this/directory/should/not/exist"},
        func() gitCommandRunner {
          return &mockGitCommandRunner{}
        }(),
      )

      a := assert.New(t)
      a.Error(err)
      a.Equal(errFilesystem, err.code)
    },
    t,
  )
}
