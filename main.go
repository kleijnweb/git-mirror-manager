package main

import (
	"github.com/kleijnweb/git-mirror-manager/manager"
)

func main() {
	manager.NewManagerServer().Start()
}
