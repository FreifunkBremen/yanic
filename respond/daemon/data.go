package respondd

import (
	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateData() {
	for iface, iData := range d.dataByInterface {
		if iData.Nodeinfo == nil {
			iData.Nodeinfo = &data.Nodeinfo{}
		}
		if iData.Statistics == nil {
			iData.Statistics = &data.Statistics{}
		}
		if iData.Neighbours == nil {
			iData.Neighbours = &data.Neighbours{}
		}
		d.updateNodeinfo(iface, iData)
		d.updateStatistics(iface, iData)
		d.updateNeighbours(iface, iData)
	}
}

func (d *Daemon) getData(iface string) *data.ResponseData {
	d.dataMX.Lock()
	defer d.dataMX.Unlock()
	if iData, ok := d.dataByInterface[iface]; ok {
		return iData
	}
	d.dataByInterface[iface] = &data.ResponseData{}
	d.updateData()
	return d.dataByInterface[iface]
}
