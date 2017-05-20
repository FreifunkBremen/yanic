package all

import "github.com/FreifunkBremen/yanic/runtime"

func (f filterConfig) Blacklist() filterFunc {
	v, ok := f["blacklist"]
	if !ok {
		return noFilter
	}

	list := make(map[string]interface{})
	for _, nodeid := range v.([]interface{}) {
		list[nodeid.(string)] = true
	}

	return func(node *runtime.Node) *runtime.Node {
		if nodeinfo := node.Nodeinfo; nodeinfo != nil {
			if _, ok := list[nodeinfo.NodeID]; ok {
				return nil
			}
		}
		return node
	}
}
