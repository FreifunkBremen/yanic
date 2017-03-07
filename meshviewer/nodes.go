package meshviewer

import (
	"log"
	"time"

	"github.com/FreifunkBremen/yanic/jsontime"
	"github.com/FreifunkBremen/yanic/runtime"
)

// NodesV1 struct, to support legacy meshviewer (which are in master branch)
//  i.e. https://github.com/ffnord/meshviewer/tree/master
type NodesV1 struct {
	Version   int              `json:"version"`
	Timestamp jsontime.Time    `json:"timestamp"` // Timestamp of the generation
	List      map[string]*Node `json:"nodes"`     // the current nodemap, indexed by node ID
}

// NodesV2 struct, to support new version of meshviewer (which are in legacy develop branch or newer)
//  i.e. https://github.com/ffnord/meshviewer/tree/dev or https://github.com/ffrgb/meshviewer/tree/develop
type NodesV2 struct {
	Version   int           `json:"version"`
	Timestamp jsontime.Time `json:"timestamp"` // Timestamp of the generation
	List      []*Node       `json:"nodes"`     // the current nodemap, as array
}

// GetNodesV1 transform data to legacy meshviewer
func GetNodesV1(nodes *runtime.Nodes) *NodesV1 {
	meshviewerNodes := &NodesV1{
		Version:   1,
		List:      make(map[string]*Node),
		Timestamp: jsontime.Now(),
	}

	for nodeID := range nodes.List {
		nodeOrigin := nodes.List[nodeID]

		if nodeOrigin.Statistics == nil {
			continue
		}

		node := &Node{
			Firstseen: nodeOrigin.Firstseen,
			Lastseen:  nodeOrigin.Lastseen,
			Flags: Flags{
				Online:  nodeOrigin.Online,
				Gateway: nodeOrigin.Gateway,
			},
			Nodeinfo: nodeOrigin.Nodeinfo,
		}
		node.Statistics = NewStatistics(nodeOrigin.Statistics)
		meshviewerNodes.List[nodeID] = node
	}
	return meshviewerNodes
}

// GetNodesV2 transform data to modern meshviewers
func GetNodesV2(nodes *runtime.Nodes) *NodesV2 {
	meshviewerNodes := &NodesV2{
		Version:   2,
		Timestamp: jsontime.Now(),
	}

	for nodeID := range nodes.List {
		nodeOrigin := nodes.List[nodeID]
		if nodeOrigin.Statistics == nil {
			continue
		}
		node := &Node{
			Firstseen: nodeOrigin.Firstseen,
			Lastseen:  nodeOrigin.Lastseen,
			Flags: Flags{
				Online:  nodeOrigin.Online,
				Gateway: nodeOrigin.Gateway,
			},
			Nodeinfo: nodeOrigin.Nodeinfo,
		}
		node.Statistics = NewStatistics(nodeOrigin.Statistics)
		meshviewerNodes.List = append(meshviewerNodes.List, node)
	}
	return meshviewerNodes
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
		switch version {
		case 1:
			runtime.SaveJSON(GetNodesV1(nodes), path)
		case 2:
			runtime.SaveJSON(GetNodesV2(nodes), path)
		default:
			log.Panicf("invalid nodes version: %d", version)
		}
	}

	if path := config.Meshviewer.GraphPath; path != "" {
		runtime.SaveJSON(BuildGraph(nodes), path)
	}
}
