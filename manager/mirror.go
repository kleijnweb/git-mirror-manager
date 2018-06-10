package manager

import (
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type mirror struct {
	Uri            string
	Name           string
	Path           string
	Cron           *cron.Cron
	UpdateInterval string
}

// mirrorNameFromUri Creates a name from a Git uri.
// It will panic if the URI is not in the expected format.
func mirrorNameFromUri(uri string) (name string) {
	parts := strings.Split(uri, "/")
	name = parts[len(parts)-2]
	name += "/" + parts[len(parts)-1]
	name = strings.ToLower(strings.TrimSuffix(name, ".git"))
	return
}

func (m *mirror) init(config *config) *Error {
	if m.Uri == "" {
		return newError("mirror uri cannot be empty", errUser)
	}

	m.Name = mirrorNameFromUri(m.Uri)
	m.UpdateInterval = config.mirrorUpdateInterval
	m.Path = config.mirrorBaseDir + "/" + m.Name

	log.Infof("Expecting repository at '%m'", m.Path)

	if _, err := os.Stat(m.Path); err != nil {
		if os.IsNotExist(err) {
			if err := m.assertValidRemote(m.Uri); err != nil {
				return err
			}
			go func() {
				if err := m.clone(); err != nil {
					log.Error(err)
				}
			}()
		} else {
			return &Error{err, errFilesystem}
		}
	}

	m.createCron()

	log.Printf("Initialized mirror '%s'", m.Name)

	return nil
}

func (m *mirror) destroy() *Error {
	m.Cron.Stop()
	return m.removeData()
}

func (m *mirror) update() *Error {
	log.Printf("Updating '%s'", m.Name)
	if _, err := gitFetchPrune(m.Path); err != nil {
		return err
	}

	log.Printf("Updating '%s' completed", m.Name)
	return nil
}

func (m *mirror) assertValidRemote(uri string) *Error {
	log.Printf("Testing '%s'", uri)
	if _, err := gitLsRemoteTags(uri); err != nil {
		return err
	}
	log.Info("Test passed")
	return nil
}

func (m *mirror) clone() *Error {
	log.Infof("Cloning '%s'", m.Name)
	_, err := gitCreateMirror(m.Uri, m.Path)
	log.Infof("Cloning '%s' completed", m.Name)
	return err
}

func (m *mirror) createDists() *Error {
	output, err := gitLsRemoteTags(m.Uri)
	if err != nil {
		return err
	}
	for _, tag := range strings.Split("\n", string(output)) {
		_, err := gitCreateTagArchive(tag, m.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mirror) removeData() *Error {
	log.Infof("Removing directory '%s'", m.Path)
	if err := os.RemoveAll(m.Path); err != nil {
		return &Error{err, errFilesystem}
	}
	log.Infof("Done removing '%s'", m.Path)
	return nil
}

func (m *mirror) createCron() *Error {
	if strings.ToLower(m.UpdateInterval) == "false" {
		return nil
	}
	m.Cron = cron.New()
	if err := m.Cron.AddFunc(m.UpdateInterval, func() {
		if err := m.update(); err != nil {
			err.Log()
		}
	}); err != nil {
		return &Error{err, errCron}
	}

	m.Cron.Start()
	return nil
}
