package respondd

import (
	"net"
	"time"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/data"
)

func (d *Daemon) Start() {
	d.dataByInterface = make(map[string]*data.ResponseData)
	if d.AnswerByZones == nil {
		d.AnswerByZones = make(map[string]*AnswerConfig)
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
