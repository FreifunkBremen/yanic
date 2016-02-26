package respond

import "time"

//Daemon struct
type Daemon struct {
	collectors []*Collector
}

//NewDaemon create a list of collectors
func NewDaemon(parseFunc func(coll *Collector, res *Response)) *Daemon {
	return &Daemon{
		collectors: []*Collector{
			NewCollector("statistics", parseFunc),
			NewCollector("nodeinfo", parseFunc),
			NewCollector("neighbours", parseFunc),
		},
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
