package respondd

import (
	"errors"
	"net"

	"github.com/Vivena/babelweb2/parser"

	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateNeighbours(iface string, resp *data.ResponseData) {
	_, nodeID := d.getAnswer(iface)
	resp.Neighbours.NodeID = nodeID
	resp.Neighbours.Batadv = make(map[string]data.BatadvNeighbours)

	for _, bface := range d.Batman {

		b := NewBatman(bface)
		if b == nil {
			continue
		}

		for bfaceAddr, n := range b.Neighbours() {
			resp.Neighbours.Batadv[bfaceAddr] = n
		}
	}

	if d.babelData == nil {
		return
	}

	resp.Neighbours.Babel = make(map[string]data.BabelNeighbours)
	d.babelData.Iter(func(t parser.Transition) error {
		if t.Table != "interface" {
			return nil
		}
		if t.Data["up"].(bool) {
			addr := t.Data["ipv6"].(net.IP)
			resp.Neighbours.Babel[string(t.Field)] = data.BabelNeighbours{
				Protocol:         "babel",
				LinkLocalAddress: addr.String(),
				Neighbours:       make(map[string]data.BabelLink),
			}
		}
		return nil
	})

	d.babelData.Iter(func(t parser.Transition) error {
		if t.Table != "neighbour" {
			return nil
		}
		ifname, ok := t.Data["if"].(string)
		if !ok {
			return errors.New("neighbour without if")
		}
		addr, ok := t.Data["address"].(net.IP)
		if !ok {
			return errors.New("neighbour without address")
		}
		if bIfname, ok := resp.Neighbours.Babel[ifname]; ok {
			link := data.BabelLink{
				RXCost: int(t.Data["rxcost"].(uint64)),
				TXCost: int(t.Data["txcost"].(uint64)),
				Cost:   int(t.Data["cost"].(uint64)),
			}
			bIfname.Neighbours[addr.String()] = link
			return nil
		}
		return errors.New("ifname not found during parsing")
	})
}
