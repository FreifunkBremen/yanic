package respondd

import (
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateStatistics(iface string, resp *data.ResponseData) {
	config, nodeID := d.getAnswer(iface)

	resp.Statistics.NodeID = nodeID

	if uptime, err := host.Uptime(); err == nil {
		resp.Statistics.Uptime = float64(uptime)
	}
	if m, err := mem.VirtualMemory(); err == nil {
		resp.Statistics.Memory.Cached = int64(m.Cached)
		resp.Statistics.Memory.Total = int64(m.Total)
		resp.Statistics.Memory.Buffers = int64(m.Buffers)
		resp.Statistics.Memory.Free = int64(m.Free)
		resp.Statistics.Memory.Available = int64(m.Available)
	}
	if v, err := load.Avg(); err == nil {
		resp.Statistics.LoadAverage = v.Load1
	}
	if v, err := load.Misc(); err == nil {
		resp.Statistics.Processes.Running = uint32(v.ProcsRunning)
		resp.Statistics.Processes.Total = uint32(v.ProcsTotal)
	}
	if ls, err := net.IOCounters(true); err == nil {
		resp.Statistics.Traffic.Tx = &data.Traffic{}
		resp.Statistics.Traffic.Rx = &data.Traffic{}

		allowedInterfaces := make(map[string]bool)

		for _, iface := range config.InterfacesTraffic {
			allowedInterfaces[iface] = true
		}

		for _, s := range ls {
			if i, ok := allowedInterfaces[s.Name]; !ok || !i {
				continue
			}
			resp.Statistics.Traffic.Tx.Bytes = float64(s.BytesSent)
			resp.Statistics.Traffic.Tx.Packets = float64(s.PacketsSent)
			resp.Statistics.Traffic.Tx.Dropped = float64(s.Dropout)

			resp.Statistics.Traffic.Rx.Bytes = float64(s.BytesRecv)
			resp.Statistics.Traffic.Rx.Packets = float64(s.PacketsRecv)
			resp.Statistics.Traffic.Rx.Dropped = float64(s.Dropin)
		}
	}
}
