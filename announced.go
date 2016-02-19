package main
import (
	"time"
)

type Announced struct {
	NodeServer	*NodeServer
	nodes           *Nodes
	outputFile      string
	collectInterval time.Duration
	saveInterval    time.Duration
	collectors	[]*Collector
}

func NewAnnounced(ns NodeServer) *Announced {
	collects := []*Collector{
		NewCollector("statistics"),
		NewCollector("nodeinfo"),
		NewCollector("neighbours"),
	}
	return &Announced{
		ns,
		NewNodes(),
		output,
		time.Second * time.Duration(15),
		time.Second * time.Duration(15),
		collects,
	}
}
func (announced *Announced) Run() {

}
