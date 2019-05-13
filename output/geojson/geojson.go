package geojson

import (
	"github.com/FreifunkBremen/yanic/lib/jsontime"
	"github.com/FreifunkBremen/yanic/runtime"
)

type GeoJSON struct {
	Type      string        `json:"type"`
	Timestamp jsontime.Time `json:"updated_at"` // Timestamp of the generation
	Features  []*Feature    `json:"features"`
}

type Feature struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
	Geometry   Geometry          `json:"geometry"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func NewPoint(n *runtime.Node) *Feature {
	if n.Nodeinfo == nil || n.Nodeinfo.Location == nil {
		return nil
	}
	location := n.Nodeinfo.Location

	return &Feature{
		Type: "Feature",
		Properties: map[string]string{
			"name": n.Nodeinfo.Hostname,
		},
		Geometry: Geometry{
			Type: "Point",
			Coordinates: []float64{
				location.Longitude,
				location.Latitude,
			},
		},
	}
}

func transform(nodes *runtime.Nodes) *GeoJSON {
	nodelist := &GeoJSON{
		Type:      "FeatureCollection",
		Timestamp: jsontime.Now(),
	}

	for _, nodeOrigin := range nodes.List {
		point := NewPoint(nodeOrigin)
		if point != nil {
			nodelist.Features = append(nodelist.Features, point)
		}
	}
	return nodelist
}
