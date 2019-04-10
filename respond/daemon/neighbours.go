package respondd

import (
	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateNeighbours(iface string, resp *data.ResponseData) {
	_, nodeID := d.getAnswer(iface)
	resp.Neighbours.NodeID = nodeID
	resp.Neighbours.Batadv = make(map[string]data.BatadvNeighbours)
	for _, bface := range d.Batman {
		b := NewBatman(bface)
		for bfaceAddr, n := range b.Neighbours() {
			resp.Neighbours.Batadv[bfaceAddr] = n
		}
	}
}
