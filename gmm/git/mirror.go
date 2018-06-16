package git

import (
  "github.com/kleijnweb/git-mirror-manager/gmm"
  "github.com/kleijnweb/git-mirror-manager/gmm/util"
  "github.com/robfig/cron"
  log "github.com/sirupsen/logrus"
  "os"
  "strings"
)

type Mirror struct {
  Name           string
  uri            string
  path           string
  cron           *cron.Cron
  updateInterval string
  cmd            CommandRunner
  fs             util.FileSystemUtil
}

func NewMirror(uri string, baseDir string, updateInterval string, cmd CommandRunner, fs util.FileSystemUtil, ) (*Mirror, gmm.ApplicationError) {

  if uri == "" {
    return nil, gmm.NewError("mirror uri cannot be empty", gmm.ErrUser)
  }

  name := MirrorNameFromURI(uri)

  m := &Mirror{
    Name:           name,
    updateInterval: updateInterval,
    uri:            uri,
    path:           baseDir + "/" + name,
    cmd:            cmd,
    fs:             fs,
  }

  log.Infof("Expecting repository at '%m'", m.path)

  if ! m.fs.DirectoryExists(m.path) {
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

  m.createCron()

  log.Printf("Initialized mirror '%s'", m.Name)

  return m, nil
}

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

func (m *Mirror) Destroy() gmm.ApplicationError {
  m.cron.Stop()
  return m.removeData()
}

func (m *Mirror) UpdateInterval() string {
  return m.updateInterval
}

func (m *Mirror) Path() string {
  return m.path
}

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

func (m *Mirror) createCron() gmm.ApplicationError {
  if strings.ToLower(m.updateInterval) == "false" {
    return nil
  }
  m.cron = cron.New()
  if err := m.cron.AddFunc(m.updateInterval, func() {
    if err := m.Update(); err != nil {
      log.Error(err)
    }
  }); err != nil {
    return gmm.NewErrorUsingError(err, gmm.ErrCron)
  }

  m.cron.Start()
  return nil
}
