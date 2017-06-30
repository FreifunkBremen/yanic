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
func BuildNodesV2(toFilter filter, nodes *runtime.Nodes) interface{} {
	meshviewerNodes := &NodesV2{
		Version:   2,
		Timestamp: jsontime.Now(),
	}

	for _, nodeOrigin := range nodes.List {
		nodeFiltere := toFilter(nodeOrigin)
		if nodeOrigin.Statistics == nil || nodeFiltere == nil {
			continue
		}
		node := &Node{
			Firstseen: nodeFiltere.Firstseen,
			Lastseen:  nodeFiltere.Lastseen,
			Flags: Flags{
				Online:  nodeFiltere.Online,
				Gateway: nodeFiltere.IsGateway(),
			},
			Nodeinfo: nodeFiltere.Nodeinfo,
		}
		node.Statistics = NewStatistics(nodeFiltere.Statistics)
		meshviewerNodes.List = append(meshviewerNodes.List, node)
	}
	return meshviewerNodes
}
