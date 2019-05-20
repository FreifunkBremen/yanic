package geojson

import (
	"strconv"
	"strings"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/paulmach/go.geojson"
)

const (
	POINT_UMAP_CLASS         = "Circle"
	POINT_UMAP_ONLINE_COLOR  = "Green"
	POINT_UMAP_OFFLINE_COLOR = "Red"
)

func newNodePoint(n *runtime.Node) (point *geojson.Feature) {
	nodeinfo := n.Nodeinfo
	location := nodeinfo.Location
	point = geojson.NewPointFeature([]float64{
		location.Longitude,
		location.Latitude,
	})
	point.Properties["id"] = nodeinfo.NodeID
	point.Properties["name"] = nodeinfo.Hostname

	point.Properties["online"] = n.Online
	var description strings.Builder
	if n.Online {
		description.WriteString("Online\n")
		if statistics := n.Statistics; statistics != nil {
			point.Properties["clients"] = statistics.Clients.Total
			description.WriteString("Clients: " + strconv.Itoa(int(statistics.Clients.Total)) + "\n")
		}
	} else {
		description.WriteString("Offline\n")
	}
	if nodeinfo.Hardware.Model != "" {
		point.Properties["model"] = nodeinfo.Hardware.Model
		description.WriteString("Model: " + nodeinfo.Hardware.Model + "\n")
	}
	if fw := nodeinfo.Software.Firmware; fw.Release != "" {
		point.Properties["firmware"] = fw.Release
		description.WriteString("Firmware: " + fw.Release + "\n")
	}
	if nodeinfo.System.SiteCode != "" {
		point.Properties["site"] = nodeinfo.System.SiteCode
		description.WriteString("Site: " + nodeinfo.System.SiteCode + "\n")
	}
	if nodeinfo.System.DomainCode != "" {
		point.Properties["domain"] = nodeinfo.System.DomainCode
		description.WriteString("Domain: " + nodeinfo.System.DomainCode + "\n")
	}
	if owner := nodeinfo.Owner; owner != nil && owner.Contact != "" {
		point.Properties["contact"] = owner.Contact
		description.WriteString("Contact: " + owner.Contact + "\n")
	}

	point.Properties["description"] = description.String()
	point.Properties["_umap_options"] = getUMapOptions(n)
	return
}

func getUMapOptions(n *runtime.Node) map[string]string {
	result := map[string]string{
		"iconClass": POINT_UMAP_CLASS,
	}
	if n.Online {
		result["color"] = POINT_UMAP_ONLINE_COLOR
	} else {
		result["color"] = POINT_UMAP_OFFLINE_COLOR
	}
	return result
}

func transform(nodes *runtime.Nodes) *geojson.FeatureCollection {
	nodelist := geojson.NewFeatureCollection()

	for _, n := range nodes.List {
		if n.Nodeinfo == nil || n.Nodeinfo.Location == nil {
			continue
		}
		point := newNodePoint(n)
		if point != nil {
			nodelist.Features = append(nodelist.Features, point)
		}
	}
	return nodelist
}
