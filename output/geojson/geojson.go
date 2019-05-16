package geojson

import (
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/paulmach/go.geojson"
)

func transform(nodes *runtime.Nodes) *geojson.FeatureCollection {
	nodelist := geojson.NewFeatureCollection()

	for _, n := range nodes.List {
		if n.Nodeinfo == nil || n.Nodeinfo.Location == nil {
			continue
		}
		nodeinfo := n.Nodeinfo
		location := nodeinfo.Location
		point := geojson.NewPointFeature([]float64{
			location.Longitude,
			location.Latitude,
		})
		point.Properties["id"] = nodeinfo.NodeID
		point.Properties["name"] = nodeinfo.Hostname

		nodelist.Features = append(nodelist.Features, point)
	}
	return nodelist
}
