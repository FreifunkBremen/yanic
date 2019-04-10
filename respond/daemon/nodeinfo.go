package respondd

import (
	"fmt"
	"os"
	"runtime"
	"net"

	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) updateNodeinfo(iface string, data *data.ResponseData) {
	config, nodeID := d.getAnswer(iface)
	data.Nodeinfo.NodeID = nodeID

	if config.Hostname == "" {
		data.Nodeinfo.Hostname, _ = os.Hostname()
	} else {
		data.Nodeinfo.Hostname = config.Hostname
	}

	data.Nodeinfo.VPN = config.VPN
	data.Nodeinfo.Location = config.Location

	data.Nodeinfo.System.SiteCode = config.SiteCode
	data.Nodeinfo.System.DomainCode = config.DomainCode

	data.Nodeinfo.Hardware.Nproc = runtime.NumCPU()

	if data.Nodeinfo.Network.Mac == "" {
		data.Nodeinfo.Network.Mac = fmt.Sprintf("%s:%s:%s:%s:%s:%s", nodeID[0:2], nodeID[2:4], nodeID[4:6], nodeID[6:8], nodeID[8:10], nodeID[10:12])
	}

	if iface == "" {
		data.Nodeinfo.Network.Addresses = []string{}
		for _, i := range d.Interfaces {
			addrs := getAddresses(i)
			data.Nodeinfo.Network.Addresses = append(data.Nodeinfo.Network.Addresses, addrs...)
		}
	} else {
		data.Nodeinfo.Network.Addresses = getAddresses(iface)
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
