package meshviewer

import (
	"log"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
)

type nodeBuilder func(*runtime.Nodes) interface{}

var nodeFormats = map[int]nodeBuilder{
	1: BuildNodesV1,
	2: BuildNodesV2,
}

// Start all services to manage Nodes
func Start(config *runtime.Config, nodes *runtime.Nodes) {
	go worker(config, nodes)
}

// Periodically saves the cached DB to json file
func worker(config *runtime.Config, nodes *runtime.Nodes) {
	c := time.Tick(config.Nodes.SaveInterval.Duration)

	for range c {
		saveMeshviewer(config, nodes)
	}
}

func saveMeshviewer(config *runtime.Config, nodes *runtime.Nodes) {
	// Locking foo
	nodes.RLock()
	defer nodes.RUnlock()
	if path := config.Meshviewer.NodesPath; path != "" {
		version := config.Meshviewer.Version
		builder := nodeFormats[version]

		if builder != nil {
			runtime.SaveJSON(builder(nodes), path)
		} else {
			log.Panicf("invalid nodes version: %d", version)
		}

	}

	if path := config.Meshviewer.GraphPath; path != "" {
		runtime.SaveJSON(BuildGraph(nodes), path)
	}
}
