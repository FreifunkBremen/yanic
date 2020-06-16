package blocklist

import (
	"errors"

	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
)

type blocklist map[string]interface{}

func init() {
	filter.Register("blocklist", build)
}

func build(config interface{}) (filter.Filter, error) {
	values, ok := config.([]interface{})
	if !ok {
		return nil, errors.New("invalid configuration, array (of strings) expected")
	}

	list := make(blocklist)
	for _, value := range values {
		if nodeid, ok := value.(string); ok {
			list[nodeid] = struct{}{}
		} else {
			return nil, errors.New("invalid configuration, array of strings expected")
		}
	}
	return &list, nil
}

func (list blocklist) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		if _, ok := list[nodeinfo.NodeID]; ok {
			return nil
		}
	}
	return node
}
