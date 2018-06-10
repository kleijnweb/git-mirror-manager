package manager

import (
	"testing"
	"reflect"
	"os"
)
var defaultsTests = []struct {
	field        string
	defaultValue string
	customValue  string
	envKey       string
}{
	{"distDir", "/opt/data/dist", "/opt/data/distSomethingElse", "GIT_MIRROR_DISTDIR"},
	{"mirrorBaseDir", "/opt/data/mirrors", "/opt/data/mirrorsSomethingElse", "GIT_MIRROR_BASEDIR"},
	{"mirrorUpdateInterval", "0 * * * *", "5 * * * *", "GIT_MIRROR_UPDATE_INTERVAL"},
	{"managerAddr", ":8080", ":555", "GIT_MIRROR_MANAGER_ADDR"},
}

func TestNewConfigReadsEnv(t *testing.T) {

	for _, tt := range defaultsTests {
		t.Run(tt.field, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.customValue)
			config := newConfig()
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
	config := newConfig()
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
