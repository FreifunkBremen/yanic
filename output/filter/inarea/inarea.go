package inarea

import (
	"errors"

	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
)

type area struct {
	latitudeMin  float64
	latitudeMax  float64
	longitudeMin float64
	longitudeMax float64
}

func init() {
	filter.Register("inarea", build)
}

func build(config interface{}) (filter.Filter, error) {

	m, ok := config.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid configuration for inarea filter")
	}

	a := area{}
	if v, ok := m["latitude_min"]; ok {
		a.latitudeMin = v.(float64)
	}
	if v, ok := m["latitude_max"]; ok {
		a.latitudeMax = v.(float64)
	}
	if v, ok := m["longitude_min"]; ok {
		a.longitudeMin = v.(float64)
	}
	if v, ok := m["longitude_max"]; ok {
		a.longitudeMax = v.(float64)
	}

	// TODO bessere Fehlerbehandlung!

	return &a, nil
}

func (a *area) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		location := nodeinfo.Location
		if location == nil {
			return node
		}
		if location.Latitude >= a.latitudeMin && location.Latitude <= a.latitudeMax && location.Longitude >= a.longitudeMin && location.Longitude <= a.longitudeMax {
			return node
		}
	}
	return nil
}
