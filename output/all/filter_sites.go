package all

import "github.com/FreifunkBremen/yanic/runtime"

func (f filterConfig) Sites() filterFunc {
	v, ok := f["sites"]
	if !ok {
		return noFilter
	}

	list := make(map[string]interface{})
	for _, site := range v.([]interface{}) {
		list[site.(string)] = true
	}

	return func(node *runtime.Node) *runtime.Node {
		if nodeinfo := node.Nodeinfo; nodeinfo != nil {
			if _, ok := list[nodeinfo.System.SiteCode]; ok {
				return node
			}
		}
		return nil
	}
}
