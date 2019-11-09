package nodelist

import (
	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/lib/jsontime"
	"github.com/FreifunkBremen/yanic/runtime"
)

// Node struct
type RawNode struct {
	Firstseen  jsontime.Time    `json:"firstseen"`
	Lastseen   jsontime.Time    `json:"lastseen"`
	Online     bool             `json:"online"`
	Statistics *data.Statistics `json:"statistics"`
	Nodeinfo   *data.Nodeinfo   `json:"nodeinfo"`
	Neighbours *data.Neighbours `json:"neighbours"`
}

type NodeList struct {
	Version   string        `json:"version"`
	Timestamp jsontime.Time `json:"updated_at"` // Timestamp of the generation
	List      []*RawNode    `json:"nodes"`
}

func transform(nodes *runtime.Nodes) *NodeList {
	nodelist := &NodeList{
		Version:   "1.0.0",
		Timestamp: jsontime.Now(),
	}

	for _, nodeOrigin := range nodes.List {
		if nodeOrigin != nil {
			node := &RawNode{
				Firstseen:  nodeOrigin.Firstseen,
				Lastseen:   nodeOrigin.Lastseen,
				Online:     nodeOrigin.Online,
				Statistics: nodeOrigin.Statistics,
				Nodeinfo:   nodeOrigin.Nodeinfo,
				Neighbours: nodeOrigin.Neighbours,
			}
			nodelist.List = append(nodelist.List, node)
		}
	}
	return nodelist
}
