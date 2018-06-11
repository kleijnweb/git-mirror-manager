package manager

import (
  "github.com/kleijnweb/git-mirror-manager/internal/git"
  "github.com/kleijnweb/git-mirror-manager/internal/util"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
)

type Manager struct {
  mirrorFactory func(uri string) (*git.Mirror, *util.ApplicationError)
  mirrors       map[string]*git.Mirror
  config        *Config
  git           git.CommandRunner
  fs            util.FileSystemUtil
}

func NewManager(config *Config, mirrorFactory func(uri string) (*git.Mirror, *util.ApplicationError), git git.CommandRunner, fs util.FileSystemUtil) *Manager {
  return &Manager{
    mirrorFactory: mirrorFactory,
    config:        config,
    mirrors:       make(map[string]*git.Mirror),
    git:           git,
    fs:            fs,
  }
}

func (m *Manager) Has(name string) bool {
  _, ok := m.mirrors[name]
  return ok
}

func (m *Manager) Add(uri string) *util.ApplicationError {
  name := git.MirrorNameFromURI(uri)

  if m.Has(name) {
    return util.NewError("mirror '"+name+"' already exists", util.ErrUser)
  }

  if err := m.Set(uri); err != nil {
    return err
  }

  return nil
}

func (m *Manager) Set(uri string) *util.ApplicationError {

  if mirror, newErr := m.mirrorFactory(uri); newErr != nil {
    return newErr
  } else {
    m.mirrors[mirror.Name] = mirror
    log.Printf("Set remote '%s' using alias '%s'", uri, mirror.Name)
  }

  return nil
}

func (m *Manager) Remove(name string) *util.ApplicationError {
  if !m.Has(name) {
    return util.NewError("mirror '"+name+"' does not exist", util.ErrNotFound)
  }
  log.Printf("Removing '%s'", name)

  if err := m.mirrors[name].Destroy(); err != nil {
    return err
  }

  delete(m.mirrors, name)

  return nil
}

func (m *Manager) Update(name string) *util.ApplicationError {
  if !m.Has(name) {
    return util.NewError("mirror '"+name+"' does not exist", util.ErrUser)
  }
  log.Printf("Updating '%s'", name)
  return nil
}

// Load existing mirrors from disk
func (m *Manager) loadFromDisk(config *Config) *util.ApplicationError {

  namespaceDirs, err := ioutil.ReadDir(config.MirrorBaseDir)

  if err != nil {
    return &util.ApplicationError{err, util.ErrFilesystem}
  }

  for _, nf := range namespaceDirs {
    nsName := nf.Name()
    nsPath := config.MirrorBaseDir + "/" + nsName
    log.Printf("Handling namespace '%s'", nsName)
    repoDirs, err := ioutil.ReadDir(nsPath)
    if err != nil {
      return &util.ApplicationError{err, util.ErrFilesystem}
    }

    for _, f := range repoDirs {
      if remote, err := m.git.GetRemote(nsPath + "/" + f.Name()); err == nil {
        if err := m.Set(remote); err != nil {
          return err
        }
      } else {
        return err
      }
    }
  }

  return nil
}
