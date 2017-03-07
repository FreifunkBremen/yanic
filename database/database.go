package database

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
)

// Connection interface to use for implementation in e.g. influxdb
type Connection interface {
	// AddNode data for a single node
	AddNode(nodeID string, node *runtime.Node)
	AddStatistics(stats *runtime.GlobalStats, time time.Time)

	DeleteNode(deleteAfter time.Duration)

	Close()
}

// Connect function with config to get DB connection interface
type Connect func(config interface{}) (Connection, error)

/*
 * for selfbinding in use of the package all
 */

var Adapters = map[string]Connect{}

func AddDatabaseType(name string, n Connect) {
	Adapters[name] = n
}
