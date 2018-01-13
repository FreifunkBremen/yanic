package blacklist

import (
	"errors"

	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
)

type blacklist map[string]interface{}

func init() {
	filter.Register("blacklist", build)
}

func build(config interface{}) (filter.Filter, error) {
	values, ok := config.([]string)
	if !ok {
		return nil, errors.New("invalid configuration, array of strings expected")
	}

	list := make(blacklist)
	for _, nodeid := range values {
		list[nodeid] = struct{}{}
	}
	return &list, nil
}

func (list blacklist) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		if _, ok := list[nodeinfo.NodeID]; ok {
			return nil
		}
	}
	return node
}
