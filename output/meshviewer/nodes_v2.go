package meshviewer

import (
	"github.com/FreifunkBremen/yanic/jsontime"
	"github.com/FreifunkBremen/yanic/runtime"
)

// NodesV2 struct, to support new version of meshviewer (which are in legacy develop branch or newer)
//  i.e. https://github.com/ffnord/meshviewer/tree/dev or https://github.com/ffrgb/meshviewer/tree/develop
type NodesV2 struct {
	Version   int           `json:"version"`
	Timestamp jsontime.Time `json:"timestamp"` // Timestamp of the generation
	List      []*Node       `json:"nodes"`     // the current nodemap, as array
}

// BuildNodesV2 transforms data to modern meshviewers
func BuildNodesV2(nodes *runtime.Nodes) interface{} {
	meshviewerNodes := &NodesV2{
		Version:   2,
		Timestamp: jsontime.Now(),
	}

	for _, nodeOrigin := range nodes.List {
		if nodeOrigin.Statistics == nil {
			continue
		}
		node := &Node{
			Firstseen: nodeOrigin.Firstseen,
			Lastseen:  nodeOrigin.Lastseen,
			Flags: Flags{
				Online:  nodeOrigin.Online,
				Gateway: nodeOrigin.IsGateway(),
			},
			Nodeinfo: nodeOrigin.Nodeinfo,
		}
		node.Statistics = NewStatistics(nodeOrigin.Statistics)
		meshviewerNodes.List = append(meshviewerNodes.List, node)
	}
	return meshviewerNodes
}
