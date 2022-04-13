package site

import (
	"errors"

	"yanic/output/filter"
	"yanic/runtime"
)

type sites map[string]interface{}

func init() {
	filter.Register("sites", build)
}

func build(config interface{}) (filter.Filter, error) {
	values, ok := config.([]interface{})
	if !ok {
		return nil, errors.New("invalid configuration, array (of strings) expected")
	}

	list := make(sites)
	for _, value := range values {
		if nodeid, ok := value.(string); ok {
			list[nodeid] = struct{}{}
		} else {
			return nil, errors.New("invalid configuration, array of strings expected")
		}
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
