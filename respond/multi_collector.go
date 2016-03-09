package respond

import "time"

//MultiCollector struct
type MultiCollector struct {
	collectors []*Collector
}

//NewMultiCollector create a list of collectors
func NewMultiCollector(parseFunc func(coll *Collector, res *Response)) *MultiCollector {
	return &MultiCollector{
		collectors: []*Collector{
			NewCollector("statistics", parseFunc),
			NewCollector("nodeinfo", parseFunc),
			NewCollector("neighbours", parseFunc),
		},
	}
}

//ListenAndSend on Collection
func (multi *MultiCollector) ListenAndSend(collectInterval time.Duration) {
	for _, col := range multi.collectors {
		col.sender(collectInterval)
	}
}

//Close all Collections
func (multi *MultiCollector) Close() {
	for _, col := range multi.collectors {
		col.Close()
	}
}
