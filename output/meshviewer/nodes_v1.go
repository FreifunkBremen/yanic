package meshviewer

import (
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

// BuildNodesV1 transforms data to legacy meshviewer
func BuildNodesV1(toFilter filter, nodes *runtime.Nodes) interface{} {
	meshviewerNodes := &NodesV1{
		Version:   1,
		List:      make(map[string]*Node),
		Timestamp: jsontime.Now(),
	}

	for nodeID, nodeOrigin := range nodes.List {
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
		meshviewerNodes.List[nodeID] = node
	}
	return meshviewerNodes
}
