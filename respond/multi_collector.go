package respond

import (
	"github.com/ffdo/node-informant/gluon-collector/data"
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
			NewCollector("statistics", interval, data.StatisticsStruct{}, onReceive),
			NewCollector("nodeinfo", interval, data.NodeInfo{}, onReceive),
			NewCollector("neighbours", interval, data.NeighbourStruct{}, onReceive),
		},
	}
}

//Close all Collections
func (multi *MultiCollector) Close() {
	for _, col := range multi.collectors {
		col.Close()
	}
}
