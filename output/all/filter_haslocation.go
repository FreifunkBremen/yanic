package all

import "github.com/FreifunkBremen/yanic/runtime"

func (f filterConfig) HasLocation() filterFunc {
	withLocation, ok := f["has_location"].(bool)
	if !ok {
		return noFilter
	}
	return func(node *runtime.Node) *runtime.Node {
		if nodeinfo := node.Nodeinfo; nodeinfo != nil {
			if withLocation {
				if location := nodeinfo.Location; location != nil {
					return node
				}
			} else {
				if location := nodeinfo.Location; location == nil {
					return node
				}
			}
		} else if !withLocation {
			return node
		}
		return nil
	}
}
