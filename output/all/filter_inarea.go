package all

import "github.com/FreifunkBremen/yanic/runtime"

type area struct {
	latitudeMin  float64
	latitudeMax  float64
	longitudeMin float64
	longitudeMax float64
}

func (f filterConfig) InArea() filterFunc {
	if areaConfigInt, ok := f["in_area"]; ok {
		areaConfig := areaConfigInt.(map[string]interface{})
		a := area{}
		if v, ok := areaConfig["latitude_min"]; ok {
			a.latitudeMin = v.(float64)
		}
		if v, ok := areaConfig["latitude_max"]; ok {
			a.latitudeMax = v.(float64)
		}
		if v, ok := areaConfig["longitude_min"]; ok {
			a.longitudeMin = v.(float64)
		}
		if v, ok := areaConfig["longitude_max"]; ok {
			a.longitudeMax = v.(float64)
		}
		return func(node *runtime.Node) *runtime.Node {
			if nodeinfo := node.Nodeinfo; nodeinfo != nil {
				location := nodeinfo.Location
				if location == nil {
					return node
				}
				if location.Latitude >= a.latitudeMin && location.Latitude <= a.latitudeMax && location.Longtitude >= a.longitudeMin && location.Longtitude <= a.longitudeMax {
					return node
				}
			}
			return nil
		}
	}
	return noFilter
}
