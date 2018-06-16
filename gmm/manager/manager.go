package manager

import (
	"github.com/kleijnweb/git-mirror-manager/gmm"
	"github.com/kleijnweb/git-mirror-manager/gmm/git"
	"github.com/kleijnweb/git-mirror-manager/gmm/util"
	log "github.com/sirupsen/logrus"
)

// Manager provides a simple interface to mirror management
type Manager struct {
	mirrorFactory func(uri string) (*git.Mirror, gmm.ApplicationError)
	mirrors       map[string]*git.Mirror
	cmd           git.CommandRunner
	fs            util.FileSystemUtil
}

// NewManager creates a new Manager struct
func NewManager(mirrorFactory func(uri string) (*git.Mirror, gmm.ApplicationError), cmd git.CommandRunner, fs util.FileSystemUtil) *Manager {
	return &Manager{
		mirrorFactory: mirrorFactory,
		mirrors:       make(map[string]*git.Mirror),
		cmd:           cmd,
		fs:            fs,
	}
}

// HasName tests whether name corresponds to a known mirror
func (m *Manager) HasName(name string) bool {
	_, ok := m.mirrors[name]
	return ok
}

// AddByURI adds a new mirror, or fails if the name was already used
func (m *Manager) AddByURI(uri string) gmm.ApplicationError {
	name := git.MirrorNameFromURI(uri)

	if m.HasName(name) {
		return gmm.NewError("mirror '"+name+"' already exists", gmm.ErrUser)
	}

	if err := m.setByURI(uri); err != nil {
		return err
	}

	return nil
}

// RemoveByName unregisters and destroys a mirror, or fails if the name is unknown
func (m *Manager) RemoveByName(name string) gmm.ApplicationError {
	if !m.HasName(name) {
		return gmm.NewError("mirror '"+name+"' does not exist", gmm.ErrNotFound)
	}
	log.Printf("Removing '%s'", name)

	if err := m.mirrors[name].Destroy(); err != nil {
		return err
	}

	delete(m.mirrors, name)

	return nil
}

// LoadFromDisk loads existing mirrors from disk
func (m *Manager) LoadFromDisk(baseDir string) gmm.ApplicationError {

	namespaceDirs, err := m.fs.ReadDir(baseDir)

	if err != nil {
		return gmm.NewErrorUsingError(err, gmm.ErrFilesystem)
	}

	for _, nf := range namespaceDirs {
		nsName := nf.Name()
		nsPath := baseDir + "/" + nsName
		log.Printf("Handling namespace '%s'", nsName)
		repoDirs, err := m.fs.ReadDir(nsPath)
		if err != nil {
			return gmm.NewErrorUsingError(err, gmm.ErrFilesystem)
		}

		for _, f := range repoDirs {
			if remote, err := m.cmd.GetRemote(nsPath + "/" + f.Name()); err == nil {
				if err := m.setByURI(remote); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func (m *Manager) setByURI(uri string) gmm.ApplicationError {

	mirror, err := m.mirrorFactory(uri)

	if err != nil {
		return err
	}

	m.mirrors[mirror.Name] = mirror
	log.Printf("Set remote '%s' using alias '%s'", uri, mirror.Name)

	return nil
}
