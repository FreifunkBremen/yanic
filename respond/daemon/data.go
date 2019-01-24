package respondd

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/FreifunkBremen/yanic/data"
)

func trim(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\n"))
}

func (d *Daemon) updateData() {
	nodeID := ""
	// Nodeinfo
	if d.Data.Nodeinfo == nil {
		d.Data.Nodeinfo = &data.Nodeinfo{}
	} else {
		nodeID = d.Data.Nodeinfo.NodeID
	}
	if d.Data.Nodeinfo.Hostname == "" {
		d.Data.Nodeinfo.Hostname, _ = os.Hostname()
	}

	// Statistics
	if d.Data.Statistics == nil {
		d.Data.Statistics = &data.Statistics{}
	} else if nodeID == "" {
		nodeID = d.Data.Statistics.NodeID
	}

	// Neighbours
	if d.Data.Neighbours == nil {
		d.Data.Neighbours = &data.Neighbours{}
	} else if nodeID == "" {
		nodeID = d.Data.Neighbours.NodeID
	}

	if nodeID == "" && !d.MultiInstance {
		if v, err := ioutil.ReadFile("/etc/machine-id"); err == nil {
			nodeID = trim(string(v))[:12]
		}
	}
	d.Data.Nodeinfo.NodeID = nodeID
	d.Data.Statistics.NodeID = nodeID
	d.Data.Neighbours.NodeID = nodeID

	for _, data := range d.dataByInterface {
		data.Nodeinfo = d.Data.Nodeinfo
	}
}

func (d *Daemon) getData(iface string) *data.ResponseData {
	if !d.MultiInstance {
		return d.Data
	}
	if data, ok := d.dataByInterface[iface]; ok {
		return data
	}
	d.dataByInterface[iface] = &data.ResponseData{}
	d.updateData()
	return d.dataByInterface[iface]
}
