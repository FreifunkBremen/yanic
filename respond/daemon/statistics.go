package respondd

import (
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"

	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateStatistics(iface string, data *data.ResponseData) {
	_, nodeID := d.getAnswer(iface)
	data.Statistics.NodeID = nodeID
	if uptime, err := host.Uptime(); err == nil {
		data.Statistics.Uptime = float64(uptime)
	}
	if m, err := mem.VirtualMemory(); err == nil {
		data.Statistics.Memory.Cached = int64(m.Cached)
		data.Statistics.Memory.Total = int64(m.Total)
		data.Statistics.Memory.Buffers = int64(m.Buffers)
		data.Statistics.Memory.Free = int64(m.Free)
		data.Statistics.Memory.Available = int64(m.Available)
	}
	if v, err := load.Avg(); err == nil {
		data.Statistics.LoadAverage = v.Load1
	}
	if v, err := load.Misc(); err == nil {
		data.Statistics.Processes.Running = uint32(v.ProcsRunning)
		//TODO fix after upstream
		data.Statistics.Processes.Total = uint32(v.Ctxt)
	}
}
