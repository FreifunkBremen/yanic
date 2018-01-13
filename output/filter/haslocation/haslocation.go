package haslocation

import (
	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
)

type haslocation struct {
	has bool
}

func init() {
	filter.Register("haslocation", build)
}

func build(config interface{}) (filter.Filter, error) {
	return &haslocation{
		has: config.(bool),
	}, nil
}

func (h *haslocation) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		if h.has {
			if location := nodeinfo.Location; location != nil {
				return node
			}
		} else {
			if location := nodeinfo.Location; location == nil {
				return node
			}
		}
	} else if !h.has {
		return node
	}
	return nil
}
