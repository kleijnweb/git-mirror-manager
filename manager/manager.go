package manager

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"io/ioutil"
	"os/exec"
)

type manager struct {
	mirrors map[string]*mirror
	config  *config
}

// Load existing mirrors from disk
func (m *manager) configure(config *config) *Error {
	m.config = config
	m.mirrors = make(map[string]*mirror)
	return m.loadFromDisk(config)
}

// Load existing mirrors from disk
func (m *manager) loadFromDisk(config *config) *Error {

	namespaceDirs, err := ioutil.ReadDir(config.mirrorBaseDir)

	if err != nil {
		return &Error{err, errFilesystem}
	}

	for _, nf := range namespaceDirs {
		nsName := nf.Name()
		nsPath := config.mirrorBaseDir + "/" + nsName
		log.Printf("Handling ns '%s'", nsName)
		repoDirs, err := ioutil.ReadDir(nsPath)
		if err != nil {
			return &Error{err, errFilesystem}
		}

		for _, f := range repoDirs {
			fullPath := nsPath + "/" + f.Name()
			cmd := exec.Command("git", "config", "--get", "remote.origin.url")
			cmd.Dir = fullPath
			output, err := cmd.CombinedOutput()
			trimmedOutput := strings.TrimSpace(string(output))
			if err != nil {
				log.Warn("Git said: " + trimmedOutput)
				return &Error{err, errGitCommand}
			}
			mirror := &mirror{
				Uri:  trimmedOutput,
				Path: fullPath,
			}
			if err := mirror.init(config); err != nil {
				return err
			}
			m.mirrors[mirror.Name] = mirror
			log.Printf("Initialized mirror '%s'", mirror.Name)
		}
	}

	return nil
}

func (m *manager) has(name string) bool {
	_, ok := m.mirrors[name]
	return ok
}

func (m *manager) add(uri string) *Error {
	name := mirrorNameFromUri(uri)

	if m.has(name) {
		return newError("mirror '"+name+"' already exists", errUser)
	}

	mirror := &mirror{
		Uri:  uri,
		Name: name,
	}

	if err := mirror.init(m.config); err != nil {
		return err
	}
	m.mirrors[mirror.Name] = mirror

	log.Printf("Added '%s' as '%s'", uri, mirror.Name)

	return nil
}

func (m *manager) remove(name string) *Error {
	if ! m.has(name) {
		return newError("mirror '"+name+"' does not exist", errNotFound)
	}
	log.Printf("Removing '%s'", name)

	if err := m.mirrors[name].destroy(); err != nil {
		return err
	}

	delete(m.mirrors, name)

	return nil
}

func (m *manager) update(name string) *Error {
	if ! m.has(name) {
		return newError("mirror '"+name+"' does not exist", errUser)
	}
	log.Printf("Updating '%s'", name)
	return nil
}
