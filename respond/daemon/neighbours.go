package respondd

import (
	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateNeighbours(iface string, data *data.ResponseData) {
	_, nodeID := d.getAnswer(iface)
	data.Neighbours.NodeID = nodeID
}
