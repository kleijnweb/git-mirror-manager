package manager_test

import (
	"github.com/kleijnweb/git-mirror-manager/gmm"
	"github.com/kleijnweb/git-mirror-manager/gmm/git"
	"github.com/kleijnweb/git-mirror-manager/gmm/manager"
	"github.com/kleijnweb/git-mirror-manager/gmm/util"
	"github.com/kleijnweb/git-mirror-manager/mocks"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type mockedFileInfo struct {
	name string
	os.FileInfo
}

func (f *mockedFileInfo) Name() string {
	return f.name
}

var mirrorFactoryCalled = false
var fsUtilMock = &mocks.FileSystemUtil{}
var gitCommandRunnerMock = &mocks.CommandRunner{}

func NewTestManager(mirrorNames ...string) *manager.Manager {
	return manager.NewManager(
		func(uri string) (*git.Mirror, gmm.ApplicationError) {
			mirrorName := mirrorNames[0]
			mirrorNames = mirrorNames[1:]
			mirrorFactoryCalled = true

			// Stubs
			cronMock := &mocks.Cron{}
			cronMock.On("Start")
			cronMock.On("Stop")
			return &git.Mirror{Name: mirrorName, Cron: cronMock}, nil
		},
		func() git.CommandRunner {
			return gitCommandRunnerMock
		}(),
		func() util.FileSystemUtil {
			return fsUtilMock
		}(),
	)
}

func TestCannotAddSameNameMoreThanOnce(t *testing.T) {
	assertions := assert.New(t)
	m := NewTestManager("ns/a", "ns/a", "ns/b", "ns/a")
	assertions.Nil(m.AddByURI("http://example.com/ns/a"))
	assertions.Error(m.AddByURI("http://example.com/ns/a"))
	assertions.Nil(m.AddByURI("http://example.com/ns/b"))
	assertions.Error(m.AddByURI("http://example.com/ns/a"))
}

func TestAddByUriInvokesMirrorFactory(t *testing.T) {
	assertions := assert.New(t)
	m := NewTestManager("ns/a")
	m.AddByURI("http://example.com/ns/a")
	assertions.True(mirrorFactoryCalled)
}

func TestCanRemoveMirror(t *testing.T) {
	assertions := assert.New(t)
	m := NewTestManager("ns/a")
	m.AddByURI("http://example.com/ns/a")
	err := m.RemoveByName("ns/a")
	assertions.Nil(err)
}

func TestCannotRemoveNonExistentMirror(t *testing.T) {
	assertions := assert.New(t)
	m := NewTestManager("ns/a")
	err := m.RemoveByName("ns/b")
	assertions.Error(err)
	assertions.Equal(err.Code(), gmm.ErrNotFound)
}

func TestCanLoadFromDisk(t *testing.T) {
	baseDir := "/mirror/basedir"
	mirrorName := "ns/a"
	m := NewTestManager(mirrorName)
	fsUtilMock.On("ReadDir", baseDir).Return([]os.FileInfo{&mockedFileInfo{name: "ns"}}, nil)
	fsUtilMock.On("ReadDir", baseDir+"/ns").Return([]os.FileInfo{&mockedFileInfo{name: "a"}}, nil)
	gitCommandRunnerMock.On("GetRemote", baseDir+"/ns/a").Return("http://example.com/ns/a", nil)
	err := m.LoadFromDisk(baseDir)
	assertions := assert.New(t)
	assertions.Nil(err)
	assertions.True(m.HasName(mirrorName))
}
