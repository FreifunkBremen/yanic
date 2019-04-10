package respondd

import (
	"github.com/shirou/gopsutil/host"

	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateStatistics(iface string, data *data.ResponseData) {
	_, nodeID := d.getAnswer(iface)
	data.Statistics.NodeID = nodeID
	if uptime, err := host.Uptime(); err == nil {
		data.Statistics.Uptime = float64(uptime)
	}
}
