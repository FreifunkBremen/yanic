package respondd

import (
	"net"
	"time"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/lib/duration"
)

type Daemon struct {
	MultiInstance bool              `toml:"multi_instance"`
	DataInterval  duration.Duration `toml:"data_interval"`
	Listen        []struct {
		Address   string `toml:"address"`
		Interface string `toml:"interface"`
		Port      int    `toml:"port"`
	} `toml:"listen"`
	Data            *data.ResponseData `toml:"data"`
	dataByInterface map[string]*data.ResponseData
}

func (d *Daemon) Start() {
	if d.Data == nil {
		d.Data = &data.ResponseData{}
	}

	d.updateData()
	go d.updateWorker()

	for _, listen := range d.Listen {
		var socket *net.UDPConn
		var err error
		addr := net.ParseIP(listen.Address)

		if addr.IsMulticast() {
			var iface *net.Interface
			if listen.Interface != "" {
				iface, err = net.InterfaceByName(listen.Interface)
				if err != nil {
					log.Fatal(err)
				}
			}
			if socket, err = net.ListenMulticastUDP("udp6", iface, &net.UDPAddr{
				IP:   addr,
				Port: listen.Port,
			}); err != nil {
				log.Fatal(err)
			}
		} else {
			if socket, err = net.ListenUDP("udp6", &net.UDPAddr{
				IP:   addr,
				Port: listen.Port,
			}); err != nil {
				log.Fatal(err)
			}
		}
		go d.handler(socket)
	}
	log.Debug("all listener started")
}

func (d *Daemon) updateWorker() {
	c := time.Tick(d.DataInterval.Duration)

	for range c {
		d.updateData()
	}
}
