package geojson

import (
	"strconv"
	"strings"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/paulmach/go.geojson"
)

func getNodeDescription(n *runtime.Node) string {
	var description strings.Builder
	if n.Online {
		description.WriteString("Online;")
	} else {
		description.WriteString("Offline;")
	}
	if statistics := n.Statistics; statistics != nil {
		description.WriteString(" " + strconv.Itoa(int(statistics.Clients.Total)) + " Clients;")
	}
	nodeinfo := n.Nodeinfo
	if nodeinfo.Hardware.Model != "" {
		description.WriteString(" Model: " + nodeinfo.Hardware.Model + ";")
	}
	if fw := nodeinfo.Software.Firmware; fw.Release != "" {
		description.WriteString(" Firmware: " + fw.Release + ";")
	}
	if nodeinfo.System.SiteCode != "" {
		description.WriteString(" Site: " + nodeinfo.System.SiteCode + ";")
	}
	if nodeinfo.System.DomainCode != "" {
		description.WriteString(" Domain: " + nodeinfo.System.DomainCode + ";")
	}
	if owner := nodeinfo.Owner; owner != nil && owner.Contact != "" {
		description.WriteString(" Contact: " + owner.Contact + ";")
	}

	return description.String()
}

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
		point.Properties["description"] = getNodeDescription(n)

		nodelist.Features = append(nodelist.Features, point)
	}
	return nodelist
}
