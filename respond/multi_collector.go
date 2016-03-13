package respond

import (
	"github.com/FreifunkBremen/RespondCollector/data"
	"time"
)

//MultiCollector struct
type MultiCollector struct {
	collectors []*Collector
}

//NewMultiCollector create a list of collectors
func NewMultiCollector(interval time.Duration, onReceive OnReceive) *MultiCollector {
	return &MultiCollector{
		collectors: []*Collector{
			NewCollector("statistics", 0, interval, data.Statistics{}, onReceive),
			NewCollector("nodeinfo", time.Second*3, interval, data.NodeInfo{}, onReceive),
			NewCollector("neighbours", time.Second*6, interval, data.Neighbours{}, onReceive),
		},
	}
}

//Close all Collections
func (multi *MultiCollector) Close() {
	for _, col := range multi.collectors {
		col.Close()
	}
}
