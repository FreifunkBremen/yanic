package haslocation

import (
	"errors"

	"yanic/output/filter"
	"yanic/runtime"
)

type haslocation struct {
	has bool
}

func init() {
	filter.Register("has_location", build)
}

func build(config interface{}) (filter.Filter, error) {
	if value, ok := config.(bool); ok {
		return &haslocation{has: value}, nil
	}
	return nil, errors.New("invalid configuration, bool expected")
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
