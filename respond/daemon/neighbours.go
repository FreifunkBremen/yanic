package respondd

import (
	"errors"

	"github.com/Vivena/babelweb2/parser"
	"github.com/bdlm/log"

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

	log.Warn(d.babelData.String())

	resp.Neighbours.Babel = make(map[string]data.BabelNeighbours)
	d.babelData.Iter(func(bu parser.BabelUpdate) error {
		sbu := bu.ToSUpdate()
		if sbu.TableId != "interface" {
			return nil
		}
		if sbu.EntryData["up"].(bool) {
			addr := sbu.EntryData["ipv6"].(string)
			resp.Neighbours.Babel[string(sbu.EntryId)] = data.BabelNeighbours{
				Protocol:         "babel",
				LinkLocalAddress: addr,
				Neighbours:       make(map[string]data.BabelLink),
			}
		}
		return nil
	})

	d.babelData.Iter(func(bu parser.BabelUpdate) error {
		sbu := bu.ToSUpdate()
		if sbu.TableId != "neighbour" {
			return nil
		}
		ifname, ok := sbu.EntryData["if"].(string)
		if !ok {
			return errors.New("neighbour without if")
		}
		addr := sbu.EntryData["address"].(string)
		if !ok {
			return errors.New("neighbour without address")
		}
		if bIfname, ok := resp.Neighbours.Babel[ifname]; ok {
			link := data.BabelLink{
				RXCost: int(sbu.EntryData["rxcost"].(uint64)),
				TXCost: int(sbu.EntryData["txcost"].(uint64)),
				Cost:   int(sbu.EntryData["cost"].(uint64)),
			}
			bIfname.Neighbours[addr] = link
			return nil
		}
		return errors.New("ifname not found during parsing")
	})
}
