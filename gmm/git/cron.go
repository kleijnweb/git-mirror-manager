package git

import (
	"github.com/kleijnweb/git-mirror-manager/gmm"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Cron represents a simplified cron job
type Cron interface {
	Stop()
	Start()
	AddFunc(spec string, cmd func()) error
}

// CronFactory creates a new Cron
type CronFactory func(mirror *Mirror, interval string) (Cron, gmm.ApplicationError)

// CreateUpdateCron creates a Cron that updates a mirror
func CreateUpdateCron(mirror *Mirror, interval string) (Cron, gmm.ApplicationError) {
	if strings.ToLower(interval) == "false" {
		return nil, nil
	}
	c := cron.New()

  updateFn := func() {
    if err := mirror.Update(); err != nil {
      log.Error(err)
    }
  }

	if err := c.AddFunc(interval, updateFn); err != nil {
		return nil, gmm.NewErrorUsingError(err, gmm.ErrCron)
	}

	return c, nil
}
