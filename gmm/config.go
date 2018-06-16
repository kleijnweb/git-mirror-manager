package gmm

import (
	"os"
)

// Config represents application configuration
type Config struct {
	MirrorBaseDir        string
	MirrorUpdateInterval string
	ManagerAddr          string
	DistDir              string
}

// NewConfig creates application config from environment variables
func NewConfig() *Config {
	envOrDefault := func(name, fallback string) string {
		val := os.Getenv(name)
		if val == "" {
			val = fallback
		}
		return val
	}
	return &Config{
		DistDir:              envOrDefault("GIT_MIRROR_DISTDIR", "/opt/data/dist"),
		MirrorBaseDir:        envOrDefault("GIT_MIRROR_BASEDIR", "/opt/data/mirrors"),
		MirrorUpdateInterval: envOrDefault("GIT_MIRROR_UPDATE_INTERVAL", "0 0 * * *"),
		ManagerAddr:          envOrDefault("GIT_MIRROR_MANAGER_ADDR", ":8080"),
	}
}
