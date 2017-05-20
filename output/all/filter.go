package all

import "github.com/FreifunkBremen/yanic/runtime"

// Config Filter
type filterConfig map[string]interface{}

type filterFunc func(*runtime.Node) *runtime.Node

func noFilter(node *runtime.Node) *runtime.Node {
	return node
}

// Create Filter
func (f filterConfig) filtering(nodesOrigin *runtime.Nodes) *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.Config{})
	filterfuncs := []filterFunc{
		f.HasLocation(),
		f.Blacklist(),
		f.InArea(),
		f.NoOwner(),
	}

	for _, nodeOrigin := range nodesOrigin.List {
		//maybe cloning of this object is better?
		node := nodeOrigin
		for _, f := range filterfuncs {
			node = f(node)
			if node == nil {
				break
			}
		}

		if node != nil {
			nodes.AddNode(node)
		}
	}
	return nodes
}
