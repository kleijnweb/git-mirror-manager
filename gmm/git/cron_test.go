package git_test

import (
  "github.com/kleijnweb/git-mirror-manager/gmm/git"
  "github.com/robfig/cron"
  "github.com/stretchr/testify/assert"
  "testing"
  "time"
)

func TestCreateUpdateCronWillReturnNilWhenUpdateIntervalStringIsQuoteUnQuoteFalse(t *testing.T) {
  c, err := git.CreateUpdateCron(&git.Mirror{}, "false")
  assert.Nil(t, err)
  assert.Nil(t, c)
}

func TestCreateUpdateCanErrorOnAddFunc(t *testing.T) {
  _, err := git.CreateUpdateCron(&git.Mirror{}, "invalid")
  assert.Error(t, err)
}

func TestCreateUpdateCanCreateCron(t *testing.T) {
  c, err := git.CreateUpdateCron(&git.Mirror{}, "0 * * * *")
  assert.Nil(t, err)
  assert.IsType(t, &cron.Cron{}, c)
}

func TestUpdateCronIsRunnable(t *testing.T) {
  gitCommandRunnerMock.On("FetchPrune", "/some/path/ns/a").Return(nil)
  c, _ := git.CreateUpdateCron(NewTestMirror("http://example.com/ns/a", "/some/path"), "* * * * *")
  c.Start()
  time.Sleep(time.Second)
  c.Stop()
}
