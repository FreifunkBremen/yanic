package respondd

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	babelParser "github.com/Vivena/babelweb2/parser"
	"github.com/bdlm/log"

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

	if out, err := exec.Command("lsb_release", "-sri").Output(); err == nil {
		f := strings.Fields(string(out))
		if len(f) == 2 {
			resp.Nodeinfo.Software.Firmware.Base = f[0]
			resp.Nodeinfo.Software.Firmware.Release = f[1]
		}
	} else {
		log.Errorf("not able to run lsb_release: %s", err)
	}

	if out, err := exec.Command("fastd", "-v").Output(); err == nil {
		resp.Nodeinfo.Software.Fastd.Enabled = true

		f := strings.Fields(string(out))
		if len(f) >= 2 {
			resp.Nodeinfo.Software.Fastd.Version = f[1]
		}
	} else {
		log.Infof("not able to run fastd: %s", err)
	}
	if v, err := ioutil.ReadFile("/sys/module/batman_adv/version"); err == nil {
		resp.Nodeinfo.Software.BatmanAdv.Version = trim(string(v))
	}
	if babel := d.babelData; babel != nil {
		resp.Nodeinfo.Software.Babeld.Version = babel.Version()
	}

	if resp.Nodeinfo.Network.Mac == "" {
		resp.Nodeinfo.Network.Mac = fmt.Sprintf("%s:%s:%s:%s:%s:%s", nodeID[0:2], nodeID[2:4], nodeID[4:6], nodeID[6:8], nodeID[8:10], nodeID[10:12])
	}

	if iface != "" {
		resp.Nodeinfo.Network.Addresses = getAddresses(iface)
	} else {
		resp.Nodeinfo.Network.Addresses = []string{}
	}
	resp.Nodeinfo.Network.Mesh = make(map[string]*data.NetworkInterface)
	for _, bface := range d.Batman {

		b := NewBatman(bface)

		if b == nil {
			continue
		}

		mesh := data.NetworkInterface{}

		for _, bbface := range b.Interfaces {
			addr := b.Address(bbface)
			if addr != "" {
				mesh.Interfaces.Tunnel = append(mesh.Interfaces.Tunnel, addr)
			}
		}

		resp.Nodeinfo.Network.Mesh[bface] = &mesh
	}

	if d.babelData == nil {
		return
	}

	meshBabel := data.NetworkInterface{}
	resp.Nodeinfo.Network.Mesh["babel"] = &meshBabel

	d.babelData.Iter(func(t babelParser.Transition) error {
		if t.Table != "interface" {
			return nil
		}
		if t.Data["up"].(bool) {
			addrIP := t.Data["ipv6"].(net.IP)
			addr := addrIP.String()
			meshBabel.Interfaces.Tunnel = append(meshBabel.Interfaces.Tunnel, addr)
			resp.Nodeinfo.Network.Addresses = append(resp.Nodeinfo.Network.Addresses, addr)
		}
		return nil
	})
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
