package manager_test

import (
  "github.com/kleijnweb/git-mirror-manager/internal/manager"
  "os"
	"reflect"
	"testing"
)

var defaultsTests = []struct {
	field        string
	defaultValue string
	customValue  string
	envKey       string
}{
	{"DistDir", "/opt/data/dist", "/opt/data/distSomethingElse", "GIT_MIRROR_DISTDIR"},
	{"MirrorBaseDir", "/opt/data/mirrors", "/opt/data/mirrorsSomethingElse", "GIT_MIRROR_BASEDIR"},
	{"MirrorUpdateInterval", "0 * * * *", "5 * * * *", "GIT_MIRROR_UPDATE_INTERVAL"},
	{"ManagerAddr", ":8080", ":555", "GIT_MIRROR_MANAGER_ADDR"},
}

func TestNewConfigReadsEnv(t *testing.T) {

	for _, tt := range defaultsTests {
		t.Run(tt.field, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.customValue)
			config := manager.NewConfig()
			st := reflect.ValueOf(config).Elem()

			v := st.FieldByName(tt.field)
			if v.String() != tt.customValue {
				t.Errorf("got %q, want %q", v.String(), tt.customValue)
			}
			// Restore defaultValue for test isolation
			os.Setenv(tt.envKey, tt.defaultValue)
		})
	}
}

func TestNewConfigDefaults(t *testing.T) {
	config := manager.NewConfig()
	st := reflect.ValueOf(config).Elem()
	for _, tt := range defaultsTests {
		t.Run(tt.field, func(t *testing.T) {
			v := st.FieldByName(tt.field)
			if v.String() != tt.defaultValue {
				t.Errorf("got %q, want %q", v.String(), tt.defaultValue)
			}
		})
	}
}
