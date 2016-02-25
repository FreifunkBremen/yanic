package responed

import "time"

//Daemon struct
type Daemon struct {
	collectors []*Collector
}

//NewDaemon create a list of collectors
func NewDaemon(parseFunc func(coll *Collector, res *Response)) *Daemon {
	collectors := []*Collector{
		NewCollector("statistics", parseFunc),
		NewCollector("nodeinfo", parseFunc),
		NewCollector("neighbours", parseFunc),
	}
	return &Daemon{
		collectors,
	}
}

//ListenAndSend on Collection
func (daemon *Daemon) ListenAndSend(collectInterval time.Duration) {
	for _, col := range daemon.collectors {
		col.sender(collectInterval)
	}
}

//Close all Collections
func (daemon *Daemon) Close() {
	for _, col := range daemon.collectors {
		col.Close()
	}
}
