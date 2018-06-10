package manager

import (
	"os"
)

type config struct {
	mirrorBaseDir        string
	mirrorUpdateInterval string
	managerAddr          string
	distDir              string
}

func newConfig() *config {
	envOrDefault := func(name, fallback string) string {
		val := os.Getenv(name)
		if val == "" {
			val = fallback
		}
		return val
	}
	return &config{
		distDir:              envOrDefault("GIT_MIRROR_DISTDIR", "/opt/data/dist"),
		mirrorBaseDir:        envOrDefault("GIT_MIRROR_BASEDIR", "/opt/data/mirrors"),
		mirrorUpdateInterval: envOrDefault("GIT_MIRROR_UPDATE_INTERVAL", "0 * * * *"),
		managerAddr:          envOrDefault("GIT_MIRROR_MANAGER_ADDR", ":8080"),
	}
}
