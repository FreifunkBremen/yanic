package output

import "github.com/FreifunkBremen/yanic/runtime"

// Output interface to use for implementation in e.g. influxdb
type Output interface {
	// InsertNode stores statistics per node
	Save(nodes *runtime.Nodes)
}

// Register function with config to get a output interface
type Register func(config map[string]interface{}) (Output, error)

// Adapters is the list of registered output adapters
var Adapters = map[string]Register{}

func RegisterAdapter(name string, n Register) {
	Adapters[name] = n
}
