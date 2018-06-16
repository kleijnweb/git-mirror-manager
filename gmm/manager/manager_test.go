package manager_test

/*
func TestNewManagerFailsIfBaseDirDoesNotExist(t *testing.T) {
  _, err := manager.NewManager(
    &manager.Config{MirrorBaseDir: "/this/directory/should/not/exist"},
    func(config *manager.Config, uri string) (*git.Mirror, *gmm.ApplicationError) {
      return &git.Mirror{}, nil
    },
    func() git.CommandRunner {
      return &mock.GitCommandRunner{}
    }(),
  )

  a := assert.New(t)
  a.Error(err)
  a.Equal(errFilesystem, err.code)
}
*/

