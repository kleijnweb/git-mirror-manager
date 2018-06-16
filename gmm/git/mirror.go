package git

import (
	"github.com/kleijnweb/git-mirror-manager/gmm"
	"github.com/kleijnweb/git-mirror-manager/gmm/util"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// Mirror represents a Git mirror
type Mirror struct {
	Name string
	Cron Cron
	uri  string
	path string
	cmd  CommandRunner
	fs   util.FileSystemUtil
}

// NewMirror creates a new Mirror struct, cloning the remote in separate subroutine
func NewMirror(
	uri string,
	baseDir string,
	updateInterval string,
	cmd CommandRunner,
	fs util.FileSystemUtil,
	updateCronFactory CronFactory,
) (*Mirror, gmm.ApplicationError) {

	if uri == "" {
		return nil, gmm.NewError("mirror uri cannot be empty", gmm.ErrUser)
	}

	name := MirrorNameFromURI(uri)

	m := &Mirror{
		Name: name,
		uri:  uri,
		path: baseDir + "/" + name,
		cmd:  cmd,
		fs:   fs,
	}

	log.Infof("Expecting repository at '%m'", m.path)

	if !m.fs.DirectoryExists(m.path) {
		if err := m.AssertValidRemote(m.uri); err != nil {
			return nil, err
		}
		log.Infof("Repository '%m' does not exists yet", m.path)
		go func() {
			if err := m.clone(); err != nil {
				log.Error(err)
			}
		}()
	}

	updateCron, err := updateCronFactory(m, updateInterval)

	if err != nil {
		return nil, err
	}

	m.Cron = updateCron
	m.Cron.Start()

	log.Printf("Initialized mirror '%s'", m.Name)

	return m, nil
}

// AssertValidRemote ensures the repository at uri can be chatted with
func (m *Mirror) AssertValidRemote(uri string) gmm.ApplicationError {
	log.Printf("Testing '%s'", uri)
	if _, err := m.cmd.LsRemoteTags(uri); err != nil {
		return err
	}
	log.Info("Test passed")
	return nil
}

// MirrorNameFromURI Creates a Name from a Git uri.
// It will panic if the uri is not in the expected format.
func MirrorNameFromURI(uri string) (name string) {
	parts := strings.Split(uri, "/")
	name = parts[len(parts)-2]
	name += "/" + parts[len(parts)-1]
	name = strings.ToLower(strings.TrimSuffix(name, ".cmd"))
	return
}

// Destroy removes local data and jobs
func (m *Mirror) Destroy() gmm.ApplicationError {
	m.Cron.Stop()
	return m.removeData()
}

// Path returns the full path to local data
func (m *Mirror) Path() string {
	return m.path
}

// Update updates the local mirror with the remote
func (m *Mirror) Update() gmm.ApplicationError {
	log.Printf("Updating '%s'", m.Name)
	if err := m.cmd.FetchPrune(m.path); err != nil {
		return err
	}

	log.Printf("Updating '%s' completed", m.Name)
	return nil
}

func (m *Mirror) clone() gmm.ApplicationError {
	log.Infof("Cloning '%s'", m.Name)
	err := m.cmd.CreateMirror(m.uri, m.path)
	log.Infof("Cloning '%s' completed", m.Name)
	return err
}

func (m *Mirror) createDists() gmm.ApplicationError {
	output, err := m.cmd.LsRemoteTags(m.uri)
	if err != nil {
		return err
	}
	for _, tag := range strings.Split("\n", string(output)) {
		err := m.cmd.CreateTagArchive(tag, m.path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Mirror) removeData() gmm.ApplicationError {
	log.Infof("Removing directory '%s'", m.path)
	if err := os.RemoveAll(m.path); err != nil {
		return gmm.NewErrorUsingError(err, gmm.ErrFilesystem)
	}
	log.Infof("Done removing '%s'", m.path)
	return nil
}
