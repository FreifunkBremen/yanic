package respondd

import (
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateNodeinfo(iface string, resp *data.ResponseData) {
	config, nodeID := d.getAnswer(iface)
	resp.Nodeinfo.NodeID = nodeID

	if config.Hostname == "" {
		resp.Nodeinfo.Hostname, _ = os.Hostname()
	} else {
		resp.Nodeinfo.Hostname = config.Hostname
	}

	resp.Nodeinfo.VPN = config.VPN
	resp.Nodeinfo.Location = config.Location

	resp.Nodeinfo.System.SiteCode = config.SiteCode
	resp.Nodeinfo.System.DomainCode = config.DomainCode

	resp.Nodeinfo.Hardware.Nproc = runtime.NumCPU()

	if resp.Nodeinfo.Network.Mac == "" {
		resp.Nodeinfo.Network.Mac = fmt.Sprintf("%s:%s:%s:%s:%s:%s", nodeID[0:2], nodeID[2:4], nodeID[4:6], nodeID[6:8], nodeID[8:10], nodeID[10:12])
	}

	if iface != "" {
		resp.Nodeinfo.Network.Addresses = getAddresses(iface)
	}

	resp.Nodeinfo.Network.Mesh = make(map[string]*data.NetworkInterface)
	for _, bface := range d.Batman {
		b := NewBatman(bface)
		mesh := data.NetworkInterface{}
		for _, bbface := range b.Interfaces {
			addr := b.Address(bbface)
			if addr != "" {
				mesh.Interfaces.Tunnel = append(mesh.Interfaces.Tunnel, addr)
			}
		}
		resp.Nodeinfo.Network.Mesh[bface] = &mesh
	}
}

func getAddresses(iface string) (addrs []string) {
	in, err := net.InterfaceByName(iface)
	if err != nil {
		return
	}
	inAddrs, err := in.Addrs()
	if err != nil {
		return
	}
	for _, a := range inAddrs {
		var ip net.IP
		switch v := a.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		default:
			continue
		}
		if ip4 := ip.To4(); ip4 == nil {
			addrs = append(addrs, ip.String())
		}
	}
	return
}
