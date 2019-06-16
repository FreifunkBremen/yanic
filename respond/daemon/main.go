package respondd

import (
	"net"
	"time"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/respond"
)

func (d *Daemon) Start() {
	d.dataByInterface = make(map[string]*data.ResponseData)
	if d.AnswerByZones == nil {
		d.AnswerByZones = make(map[string]*AnswerConfig)
	}

	d.dataMX.Lock()
	d.updateData()
	d.dataMX.Unlock()
	go d.updateWorker()

	if d.Babel != "" {
		go d.babelConnect()
	}

	for _, listen := range d.Listen {
		var socket *net.UDPConn
		var err error
		addrString := listen.Address
		if addrString == "" {
			addrString = respond.MulticastAddressDefault
		}
		port := listen.Port
		if port == 0 {
			port = respond.PortDefault
		} else if port < 0 {
			port = 0
		}
		addr := net.ParseIP(addrString)

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
				Port: port,
			}); err != nil {
				log.Fatal(err)
			}
		} else {
			if socket, err = net.ListenUDP("udp6", &net.UDPAddr{
				IP:   addr,
				Port: port,
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
		d.dataMX.Lock()
		d.updateData()
		d.dataMX.Unlock()
	}
}
