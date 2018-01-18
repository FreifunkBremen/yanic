package site

import (
	"errors"

	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
)

type sites map[string]interface{}

func init() {
	filter.Register("sites", build)
}

func build(config interface{}) (filter.Filter, error) {
	values, ok := config.([]string)
	if !ok {
		return nil, errors.New("invalid configuration, array of strings expected")
	}

	list := make(sites)
	for _, nodeid := range values {
		list[nodeid] = struct{}{}
	}
	return &list, nil
}

func (list sites) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		if _, ok := list[nodeinfo.System.SiteCode]; ok {
			return node
		}
	}
	return nil
}
