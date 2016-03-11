package respond

import "time"

//MultiCollector struct
type MultiCollector struct {
	collectors []*Collector
}

//NewMultiCollector create a list of collectors
func NewMultiCollector(interval time.Duration, parseFunc ParseFunc) *MultiCollector {
	return &MultiCollector{
		collectors: []*Collector{
			NewCollector("statistics", interval, parseFunc),
			NewCollector("nodeinfo", interval, parseFunc),
			NewCollector("neighbours", interval, parseFunc),
		},
	}
}

//Close all Collections
func (multi *MultiCollector) Close() {
	for _, col := range multi.collectors {
		col.Close()
	}
}
